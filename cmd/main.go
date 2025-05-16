package main

import (
	"fmt"
	"monkey/evaluator"
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
	"monkey/repl"
	"os"
)

func main() {
	argCount := len(os.Args)
	if argCount == 1 {
		r()
	} else if argCount == 2 {
		f, err := os.Open(os.Args[1])
		if err != nil {
			fmt.Printf("%+v\n", err)
			os.Exit(1)
		}

		l := lexer.NewFromReader(os.Args[1], f)
		p := parser.New(l)
		program := p.ParseProgram()
		obj := evaluator.Eval(program, object.NewEnvironment())

		if obj != nil {
			if obj.Type() == object.ERROR_OBJ {
				fmt.Printf("%s\n", obj.Inspect())
				os.Exit(1)
			}
			fmt.Printf("%v\n", obj.Inspect())
		}
		os.Exit(0)
	}
}

func r() {
	fmt.Printf("Monkey REPL\n")
	repl.Start(os.Stdin, os.Stdout)
}
