package vm

import (
	"fmt"
	"monkey/ast"
	"monkey/compiler"
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
	"testing"
)

func parse(input string) *ast.Program {
	l := lexer.NewFromString("vm", input)
	p := parser.New(l)
	return p.ParseProgram()
}

func testIntegerObject(expected int64, actual object.Object) error {
	result, ok := actual.(*object.Integer)
	if !ok {
		return fmt.Errorf("object is not Integer. got=%T(%+v)", actual, actual)
	}

	if result.Value != expected {
		return fmt.Errorf("object has wrong value.  got=%d, want=%d", result.Value, expected)
	}
	return nil
}

type vmTestCase struct {
	input    string
	expected interface{}
}

func runVmTests(t *testing.T, tests []vmTestCase) {
	t.Helper()

	for _, tt := range tests {
		program := parse(tt.input)
		comp := compiler.New()
		err := comp.Compile(program)
		if err != nil {
			t.Fatalf("[%s] compiler error: %+v", tt.input, err)
		}

		vm := New(comp.Bytecode())
		err = vm.Run()
		if err != nil {
			t.Fatalf("[%s] vm err: %+v", tt.input, err)
		}

		stackElem := vm.LastPoppedStackElem()
		testExpectedObject(t, tt.input, tt.expected, stackElem)

	}
}

func testExpectedObject(t *testing.T, input string, expected interface{}, actual object.Object) {
	t.Helper()
	switch expected := expected.(type) {
	case int:
		err := testIntegerObject(int64(expected), actual)
		if err != nil {
			t.Errorf("[%s] test2IntegerObject failed: %+v", input, err)
		}
	}
}

func TestIntegerArithmetic(t *testing.T) {
	tests := []vmTestCase{
		{"1", 1},
		{"2", 2},
		{"1 + 2", 3},
		{"1 - 2", -1},
		{"1 * 2", 2},
		{"4 / 2", 2},
		{"50 / 2 * 2 + 10 - 5", 55},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"5 + 2 + 10", 17},
		{"5 + 2 * 10", 25},
		{"5 * (2 + 10)", 60},
	}
	runVmTests(t, tests)
}

func TestExecutionPath(t *testing.T) {
	input := "13 + 14 + 15 + 16"
	l := lexer.NewFromString("test", input)
	p := parser.New(l)
	program := p.ParseProgram()
	programErrors := p.Errors()
	if len(programErrors) != 0 {
		t.Fatalf("Error parsing program: %+v", programErrors)
	}

	c := compiler.New()
	err := c.Compile(program)
	if err != nil {
		t.Fatalf("Error compiling program: %+v", err)
	}
	bc := c.Bytecode()
	v := New(bc)
	err = v.Run()
	if err != nil {
		t.Fatalf("Error running program: %+v", err)
	}
	fmt.Println(input, "=", v.LastPoppedStackElem().Inspect())
}
