package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

const editorConfigFile = "editor.config"

func readPreferredEditor() string {
	if file, err := os.Open(editorConfigFile); err == nil {
		defer file.Close()
		scanner := bufio.NewScanner(file)
		if scanner.Scan() {
			return scanner.Text()
		}
	}
	return ""
}

func savePreferredEditor(editor string) {
	file, err := os.Create(editorConfigFile)
	if err != nil {
		fmt.Println("Error saving editor choice:", err)
		return
	}
	defer file.Close()

	file.WriteString(editor)
}

func help() {
	fmt.Println("Problem Solver CLI")
	fmt.Println("Usage:")
	fmt.Println("  mycli [command]")
	fmt.Println("Commands:")
	fmt.Println("  solve <year> <number>  Solves the specified problem for a given year and number")
	fmt.Println("                         If the problem folder does not exist, it sets up the environment and creates template files.")
	fmt.Println("                         If the folder exists, it compiles and runs the existing solution.")
	fmt.Println("  try <year> <number>    Similar to 'solve', but always opens the 'solve.go' file in the editor")
	fmt.Println("                         regardless of whether the folder and files already exist.")
	fmt.Println("Flags:")
	fmt.Println("  --help, -h             Show help information")
}

func openFileInEditor(editor, filePath string) {
	var cmd *exec.Cmd

	switch editor {
	case "vscode":
		cmd = exec.Command("code", filePath)
	case "nvim", "emacs":
		cmd = exec.Command(editor, filePath)
	default:
		fmt.Println("Unsupported editor. Please choose from nvim, emacs, or vscode.")
		return
	}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		fmt.Println("Error opening file in editor:", err)
		return
	}
	cmd.Wait()
}

func executeSolveGo(path string) {
	cmdRun := exec.Command("go", "run", filepath.Join(path, "solve.go"), filepath.Join(path, "input.txt"))
	var runOutput bytes.Buffer
	cmdRun.Stdout = &runOutput
	cmdRun.Stderr = &runOutput

	if err := cmdRun.Run(); err != nil {
		fmt.Println("Error running solve.go:", runOutput.String())
		return
	}

	fmt.Println("Output from solve.go:\n", runOutput.String())
}

func setupProblemEnvironment(year, number string) (string, bool, error) {
	path := filepath.Join(year, number)

	wasCreated := false

	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return "", false, err
		}
		fmt.Println("Please paste the contents of input.txt and press Enter, then Ctrl+D (or Ctrl+Z on Windows):")
		scanner := bufio.NewScanner(os.Stdin)

		file, err := os.Create(filepath.Join(path, "input.txt"))
		if err != nil {
			fmt.Println("Error creating input.txt:", err)
			return "", false, err
		}
		defer file.Close()

		for scanner.Scan() {
			file.WriteString(scanner.Text() + "\n")
		}

		if err := scanner.Err(); err != nil {
			fmt.Println("Error reading input:", err)
			return "", false, err
		}

		solveFile, err := os.Create(filepath.Join(path, "solve.go"))
		if err != nil {
			fmt.Println("Error creating solve.go file:", err)
			return "", false, err
		}
		defer solveFile.Close()

		template := `package main

import (
    "bufio"
    "fmt"
    "os"
)

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: solve <path to input.txt>")
        return
    }

    inputFilePath := os.Args[1]
    file, err := os.Open(inputFilePath)
    if err != nil {
        fmt.Println("Error opening input file:", err)
        return
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    var lines []string
    for scanner.Scan() {
        lines = append(lines, scanner.Text())
    }

    if err := scanner.Err(); err != nil {
        fmt.Println("Error reading input:", err)
        return
    }

    result := solve(lines)
    fmt.Println("Result:", result)
}

func solve(lines []string) interface{} {
    // TODO: Implement the problem-solving logic here
    return nil
}
`
		if _, err := solveFile.WriteString(template); err != nil {
			fmt.Println("Error writing to solve.go:", err)
			return "", false, err
		}
		wasCreated = true
	}
	return path, wasCreated, nil
}

func solve(args []string) {
	if len(args) < 2 {
		fmt.Println("Not enough arguments for solve command. Expected format: solve <year> <number>")
		return
	}

	year, number := args[0], args[1]
	path := filepath.Join(year, number)
	path, wasCreated, err := setupProblemEnvironment(year, number)
	if err != nil {
		fmt.Println("Error setting up problem environment:", err)
		return
	}

	if wasCreated {
		solveFilePath := filepath.Join(path, "solve.go")
		editor := readPreferredEditor()
		if editor == "" {
			fmt.Print("Choose your code editor (nvim/emacs/vscode): ")
			fmt.Scan(&editor)
			savePreferredEditor(editor)
		}
		openFileInEditor(editor, solveFilePath)
	} else {
		executeSolveGo(path)
	}
}

func try(args []string) {
	if len(args) < 2 {
		fmt.Println("Not enough arguments for try command. Expected format: try <year> <number>")
		return
	}

	year, number := args[0], args[1]
	path, _, err := setupProblemEnvironment(year, number)
	if err != nil {
		fmt.Println("Error setting up problem environment:", err)
		return
	}

	solveFilePath := filepath.Join(path, "solve.go")
	editor := readPreferredEditor()
	if editor == "" {
		fmt.Print("Choose your code editor (nvim/emacs/vscode): ")
		fmt.Scan(&editor)
		savePreferredEditor(editor)
	}
	openFileInEditor(editor, solveFilePath)
}

func main() {
	helpFlag := flag.Bool("help", false, "Show help information")
	hFlag := flag.Bool("h", false, "Show help information (shorthand)")

	flag.Parse()

	if *helpFlag || *hFlag {
		help()
		return
	}

	args := flag.Args()
	if len(args) < 3 {
		fmt.Println("Invalid arguments. Use --help or -h for usage.")
		return
	}

	command := args[0]
	switch command {
	case "solve":
		solve(args[1:])
	case "try":
		try(args[1:])
	default:
		fmt.Println("Invalid command. Use --help or -h for usage.")
	}
}
