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
	switch expectedType := expected.(type) {
	case bool:
		err := testBooleanObject(input, expectedType, actual)
		if err != nil {
			t.Errorf("testBooleanObject failed: %+v", err)
		}
	case int:
		err := testIntegerObject(int64(expectedType), actual)
		if err != nil {
			t.Errorf("[%s] test2IntegerObject failed: %+v", input, err)
		}
	}
}

func testBooleanObject(input string, expected bool, actual object.Object) error {
	result, ok := actual.(*object.Boolean)
	if !ok {
		return fmt.Errorf("[%s] object is not Boolean.  got=%T(%+v)", input, actual, actual)
	}
	if result.Value != expected {
		return fmt.Errorf("[%s] object has wrong value. got=%t, want=%t", input, result.Value, expected)
	}
	return nil
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
		{"-5", -5},
		{"-10", -10},
		{"-50 + 100 + -50", 0},
		{"(5 - 10 * 2 + 15 / 3) * 2 + -10", -30},
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

func TestBooleanExpressions(t *testing.T) {
	tests := []vmTestCase{
		{"true", true},
		{"false", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 != 2", true},
		{"true == true", true},
		{"true == false", false},
		{"false == false", true},
		{"true != false", true},
		{"false != true", true},
		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false},
		{"(1 > 2) == true", false},
		{"(1 > 2) == false", true},
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
	}
	runVmTests(t, tests)
}
