# Advent of Code CLI in Go

## Introduction

This Command Line Interface (CLI) application is specifically designed for solving problems from the Advent of Code calendar using the Go programming language. It streamlines the setup and organization of problem-solving environments for each challenge based on the year and problem number. The CLI facilitates quick access to problem files and allows for executing solutions directly from the command line. While currently focused on Go, there are plans to extend support to other programming languages in future updates.

## Installation

To use this CLI, you need to have Go installed on your machine. If you don't have Go installed, you can download it from the [official Go website](https://golang.org/dl/).

1. Clone the repository to your local machine:

   ```sh
   git clone git@github.com:PiyushMishra318/advent-of-code.git
   ```

2. Navigate to the cloned directory:

   ```sh
   cd advent-of-code
   ```

3. Build the CLI tool:

   ```sh
   go build -o advent-of-code
   ```

## Usage

After building the application, you can run it using the following syntax:

```sh
./advent-of-code [command] [year] [number]
```

### Commands

- `solve <year> <number>`: Sets up the environment for a given problem if it does not exist and compiles and runs the existing solution if it does.
- `try <year> <number>`: Similar to `solve`, but always opens the `solve.go` file in a text editor for editing.
- `--help`, `-h`: Displays help information about the CLI commands.

### Examples

1. To solve a problem for the year 2023, problem number 1:

   ```sh
   ./advent-of-code solve 2023 1
   ```

2. To work on a problem (open in an editor) for the year 2023, problem number 1:

   ```sh
   ./advent-of-code try 2023 1
   ```

3. To display help information:

   ```sh
   ./advent-of-code --help
   ```

## Contributing

Contributions to this project are welcome. Please feel free to fork the repository, make changes, and submit pull requests.

## License

This project is licensed under the MIT License - see the LICENSE file for details
