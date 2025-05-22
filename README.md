# Monkey Scripting Language

Following the book [Writing an Interpreter in Go](https://a.co/d/a7Zb1Br) by Thorsten Ball, with the following modifications:
- Unicode lexer
- File evaluation `monkey <file>`
- Floating point types
- Octal and hexadecimal integer constants
- Access to environment variables
- Process Execution (with only stdout returned)

## Building Monkey

`go build -o ./bin/monkey ./cmd/main.go`
