package main

import (
	"flag"
	"fmt"
	"monkey/evaluator"
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
	"monkey/repl"
	"os"
	"runtime/pprof"
)

var cpuProfile = flag.String("cpuprofile", "", "Store cpu profile data")
var useRepl = flag.Bool("repl", false, "Start the REPL")

func main() {
	flag.Usage = usage
	flag.Parse()

	if *cpuProfile != "" {
		f, err := os.Create(*cpuProfile)
		if err != nil {
			fmt.Errorf("unable to store profile data: %+v", err)
			os.Exit(1)
		}
		defer f.Close()
		err = pprof.StartCPUProfile(f)
		if err != nil {
			fmt.Errorf("unable to store profile data: %+v", err)
			os.Exit(1)
		}
		defer pprof.StopCPUProfile()
	}

	args := flag.Args()
	if *useRepl || len(args) == 0 {
		r()
	} else {

		if len(args) != 1 {
			flag.Usage()
			os.Exit(1)
		}
		runFile(args[0])
	}
}

func runFile(fileName string) {
	f, err := os.Open(fileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}

	l := lexer.NewFromReader(fileName, f)
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

func r() {
	fmt.Printf("Monkey REPL\n")
	repl.Start(os.Stdin, os.Stdout)
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [flags] [file]\n", os.Args[0])
	fmt.Fprintln(os.Stderr, "  file\n\tA monkey script.")
	flag.PrintDefaults()
	os.Exit(1)
}
