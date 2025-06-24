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
	switch e := expected.(type) {
	case bool:
		err := testBooleanObject(input, e, actual)
		if err != nil {
			t.Errorf("testBooleanObject failed: %+v", err)
		}
	case int:
		err := testIntegerObject(int64(e), actual)
		if err != nil {
			t.Errorf("[%s] test2IntegerObject failed: %+v", input, err)
		}
	case string:
		err := testStringObject(e, actual)
		if err != nil {
			t.Errorf("testStringObject failed: %s", err)
		}
	case []int:
		array, ok := actual.(*object.Array)
		if !ok {
			t.Errorf("object not Array: %T(%+v)", actual, actual)
			return
		}
		if len(array.Elements) != len(e) {
			t.Errorf("wrong num of elements. want=%d, got=%d", len(e), len(array.Elements))
		}
		for i, el := range e {
			err := testIntegerObject(int64(el), array.Elements[i])
			if err != nil {
				t.Errorf("testIntegerObject failed: %s", err)
			}
		}
	case map[object.HashKey]int64:
		hash, ok := actual.(*object.Hash)
		if !ok {
			t.Errorf("object is not Hash.  got=%T (%+v)", actual, actual)
			return
		}
		if len(hash.Pairs) != len(e) {
			t.Errorf("Hash has wrong number of Pairs.  wat=%d, got=%d", len(e), len(hash.Pairs))
		}
		for expectedKey, expectedValue := range e {
			pair, ok := hash.Pairs[expectedKey]
			if !ok {
				t.Errorf("no pair for given key in Pairs")
			}
			err := testIntegerObject(expectedValue, pair.Value)
			if err != nil {
				t.Errorf("testIntegerObject failed: %s", err)
			}
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
		{"!(if (false) { 5; })", true},
	}
	runVmTests(t, tests)
}

func TestConditional(t *testing.T) {
	tests := []vmTestCase{
		{"if (true) { 10 }", 10},
		{"if (true) { 10 } else { 20 }", 10},
		{"if (false) { 10 } else { 20 }", 20},
		{"if (1) { 10 }", 10},
		{"if (1 < 2) { 10 }", 10},
		{"if (1 < 2) { 10 } else { 20 }", 10},
		{"if (1 > 2) { 10 } else { 20 }", 20},
		{"if (1 > 2) { 10 }", Null},
		{"if (false) { 10 }", Null},
	}

	runVmTests(t, tests)
}

func TestGlobalLetStatements(t *testing.T) {
	tests := []vmTestCase{
		{"let one = 1; one", 1},
		{"let one = 1; let two = 2; one + two", 3},
		{"let one = 1; let two = one + one; one + two", 3},
	}
	runVmTests(t, tests)
}

func testStringObject(expected string, actual object.Object) error {
	result, ok := actual.(*object.String)
	if !ok {
		return fmt.Errorf("object is not String, got=%T(%v)", actual, actual)
	}
	if result.Value != expected {
		return fmt.Errorf("object has wrong value. got=%q, want=%q", result, expected)
	}
	return nil
}

func TestStringExpressions(t *testing.T) {
	tests := []vmTestCase{
		{`"monkey"`, "monkey"},
		{`"mon" + "key"`, "monkey"},
		{`"mon" + "key" + "banana"`, "monkeybanana"},
	}
	runVmTests(t, tests)
}
func TestArrayLiterals(t *testing.T) {
	tests := []vmTestCase{
		{"[]", []int{}},
		{"[1,2,3]", []int{1, 2, 3}},
		{"[1 + 2, 3 * 4, 5 + 6]", []int{3, 12, 11}},
	}
	runVmTests(t, tests)
}

func TestHashLiterals(t *testing.T) {
	tests := []vmTestCase{
		{
			"{}", map[object.HashKey]int64{},
		},
		{
			"{1: 2, 2: 3}",
			map[object.HashKey]int64{
				(&object.Integer{Value: 1}).HashKey(): 2,
				(&object.Integer{Value: 2}).HashKey(): 3,
			},
		},
		{
			"{1 + 1: 2 * 2, 3 + 3: 4 * 4}",
			map[object.HashKey]int64{
				(&object.Integer{Value: 2}).HashKey(): 4,
				(&object.Integer{Value: 6}).HashKey(): 16,
			},
		},
	}
	runVmTests(t, tests)
}
