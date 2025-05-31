# Monkey Scripting Language

Following the book [Writing an Interpreter in Go](https://a.co/d/a7Zb1Br) and [Writing a Compiler in Go](https://a.co/d/8QmAUQn) by Thorsten Ball, with the following modifications:
- Unicode lexer
- File evaluation `monkey <file>`
- Floating point types
- Octal and hexadecimal integer constants
- Access to environment variables
- Process execution (with only stdout returned)
- \# Comments 
- Identifiers can contain numbers and underscores.

## Built-in Functions
    - len(): The length of a string or array.
    - first(): The first element of an array.
    - last(): The last element of an array.
    - rest(): All the elements of an array after the first element.
    - push(): Adds an element to the array.
    - puts(): prints a value to stdout.
    - exec(): executes a command and returns the stdout.
    - cmp(): For strings and floating point values, return -1, 0, or 1 if the first argument is less than, equal or greater than the second.

## Building Monkey

`go build -o ./bin/monkey ./cmd/main.go`
