package evaluator

import (
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
	"strings"
	"testing"
)

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"-5", -5},
		{"10", 10},
		{"-10", -10},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"5 * 2 + 10", 20},
		{"-50 + 100 + -50", 0},
		{"50 / 2 * 2 + 10", 60},
		{"3 * 3 * 3 + 10", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(tt.input, t, evaluated, tt.expected)
	}
}

func testEval(input string) object.Object {
	l := lexer.NewFromString("test", input)
	p := parser.New(l)
	program := p.ParseProgram()
	return Eval(program, object.NewEnvironment())
}

func testIntegerObject(input string, t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("[%s] object is not Integer, got= %T(%+v)", input, obj, obj)
		return false
	}

	if result.Value != expected {
		t.Errorf("[%s] object has wrong value, actual=%d expected=%d", input, result.Value, expected)
		return false
	}

	return true
}

func testIntegerArray(input string, t *testing.T, obj object.Object, expected []int64) bool {
	result, ok := obj.(*object.Array)
	if !ok {
		t.Errorf("[%s] object is not Array, got=%T(%+v)", input, obj, obj)
		return false
	}

	expectedLen := len(expected)
	actualLen := len(result.Elements)
	if expectedLen != actualLen {
		t.Errorf("[%s] array does not have exected element count.  got=%d, expected=%d", input, actualLen, expectedLen)
		return false
	}
	for i := 0; i != expectedLen; i++ {
		if result.Elements[i].Type() != object.INTEGER_OBJ {
			t.Errorf("[%s] element has wrong type at index %d, got=%T(%+v)", input, i, result.Elements[i], result.Elements[i])
			return false
		}
		intObj := result.Elements[i].(*object.Integer)

		if intObj.Value != expected[i] {
			t.Errorf("[%s] element has wrong value at index %d, got=%T(%+v)", input, i, intObj.Value, expected[i])
			return false
		}
	}

	return true
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"true != false", true},
		{"false != true", true},
		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false},
		{"(1 > 2) == true", false},
		{"(1 > 2) == false", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(tt.input, t, evaluated, tt.expected)
	}
}

func testBooleanObject(input string, t *testing.T, obj object.Object, expected bool) bool {
	result, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf("[%s] object is not Boolean. got=%T (%+v)", input, obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("[%s] object has wrong value, got=%t, want=%t", input, result.Value, expected)
		return false
	}
	return true
}

func TestBangOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(tt.input, t, evaluated, tt.expected)
	}
}

func TestIfElseExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"if (true) { 10 }", 10},
		{"if (false) { 10 }", nil},
		{"if (1) { 10 }", 10},
		{"if (1 < 2 ) { 10 }", 10},
		{"if (1 > 2) { 10 }", nil},
		{"if (1 > 2) { 10 } else { 20 }", 20},
		{"if (1 < 2) { 10 } else { 10 }", 10},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(tt.input, t, evaluated, int64(integer))
		} else {
			testNullObject(tt.input, t, evaluated)
		}
	}
}

func testNullObject(input string, t *testing.T, obj object.Object) bool {
	if obj != NULL {
		t.Errorf("[%s]: object is not NULL, got %T(%+v)", input, obj, obj)
		return false
	}
	return true
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"return 10;", 10},
		{"return 10; 9; ", 10},
		{"return 10; 9; ", 10},
		{"return 2 * 5; 9;", 10},
		{"9; return 2 * 5; 9;", 10},
		{`if (10 > 1) {
if (10 > 1) { return 10; }
return 1;
}`, 10},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(tt.input, t, evaluated, tt.expected)
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input           string
		expectedMessage string
	}{
		{"5 + true;", "type mismatch: INTEGER + BOOLEAN"},
		{"5 + true; 5", "type mismatch: INTEGER + BOOLEAN"},
		{"-true", "unknown operator: -BOOLEAN"},
		{"true + false", "type mismatch: BOOLEAN + BOOLEAN"},
		{"5; true + false; 5;", "type mismatch: BOOLEAN + BOOLEAN"},
		{"if (10 > 1) { true + false; }", "type mismatch: BOOLEAN + BOOLEAN"},
		{"if (10 > 1) {if (10 > 1) { return true + false; }}", "type mismatch: BOOLEAN + BOOLEAN"},
		{"foobar", "identifier not found: foobar"},
		{`"Hello" - "World"`, "unknown operator: STRING - STRING"},
		{`{"name": "monkey"}[fn(x){x}];`, "unusable as hash key: FUNCTION"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		errObj, ok := evaluated.(*object.Error)
		if !ok {
			t.Errorf("[%s]: no error object returnd, got=%T(%+v)", tt.input, evaluated, evaluated)
			continue
		}

		if !strings.Contains(errObj.Message, tt.expectedMessage) {
			t.Errorf("[%s]: wrong error message. expsected=%q, got=%q", tt.input, tt.expectedMessage, errObj.Message)
		}
	}
}

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let a = 5; a;", 5},
		{"let a = 5 * 5; a;", 25},
		{"let a = 5; let b = a; b;", 5},
		{"let a = 5; let b = a; let c = a + b + 5; c;", 15},
	}

	for _, tt := range tests {
		testIntegerObject(tt.input, t, testEval(tt.input), tt.expected)
	}
}

func TestFunctionObject(t *testing.T) {
	input := "fn(x) {x +2; };"
	evaluated := testEval(input)
	fn, ok := evaluated.(*object.Function)

	if !ok {
		t.Fatalf("object is not function, got=%T (%+v)", fn, fn)
	}

	if len(fn.Parameters) != 1 {
		t.Fatalf("function has wrong parameters. Parameters=%+v", fn.Parameters)
	}

	if fn.Parameters[0].String() != "x" {
		t.Fatalf("parameter is not 'x'. got=%q", fn.Parameters[0])
	}

	expectedBody := "(x + 2)"

	if fn.Body.String() != expectedBody {
		t.Fatalf("body is not %q. got=%q", expectedBody, fn.Body.String())
	}
}

func TestFunctionApplication(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let identity = fn(x) { x; }; identity(5);", 5},
		{"let identity = fn(x) { return x; }; identity(5);", 5},
		{"let double = fn(x) { x * 2; }; double(5)", 10},
		{"let add = fn(x, y) { x + y; }; add(5, 5);", 10},
		{"let add = fn(x, y) { x + y; }; add(5 + 5, add(5, 5));", 20},
		{"fn(x){x;}(5)", 5},
	}
	for _, tt := range tests {
		testIntegerObject(tt.input, t, testEval(tt.input), tt.expected)
	}
}

func TestStringLIteral(t *testing.T) {
	input := `"Hello World"`
	evaluated := testEval(input)
	str, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("object is not String. got=%T (%+v)", evaluated, evaluated)
	}

	if str.Value != "Hello World" {
		t.Errorf("String has wrong value, got=%q, expected=%q", str.Value, "Hello World")
	}
}

func TestStringConcatenation(t *testing.T) {
	input := `"Hello" + " " + "World!"`
	evaluated := testEval(input)
	str, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("object is not String. got=%T (%+v)", evaluated, evaluated)
	}

	if str.Value != "Hello World!" {
		t.Errorf("String has wrong value. got=%q, expected=%q", str.Value, "Hello World!")
	}
}

func TestBuiltinFunction(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`len("")`, 0},
		{`len("four")`, 4},
		{`len("Hello World")`, 11},
		{`len(1)`, "argument to 'len' not supported, got INTEGER"},
		{`len("one", "2")`, "wrong number of arguments. got=2, want=1"},
		{`rest(["one"], "2")`, "wrong number of arguments. got=2, want=1"},
		{`let x = [1,2,3]; push(x, 4);`, []int64{1, 2, 3, 4}},
		{`push([1,2], 4, 5);`, "wrong number of arguments. got=3, want=2"},
		{`push(1, [1,2]);`, "argument to `push` must be ARRAY, got INTEGER"},
		{`let x = [1,2,3]; rest(x);`, []int64{2, 3}},
		{"rest(1, 2, 3)", "wrong number of arguments. got=3, want=1"},
		{"rest(1)", "argument to 'rest' must be an array, got INTEGER"},
		{`let x = [1,2,3]; push(x, -5);`, []int64{1, 2, 3, -5}},
		{`let x = [1,2,3]; len(x); x;`, []int64{1, 2, 3}},
		{`let x = [1,2,3]; first(x);`, 1},
		{`first(1);`, "argument to 'first' must be an array, got INTEGER"},
		{`first([1], [2]);`, "wrong number of arguments. got=2, want=1"},
		{`last([1, 2, 3]);`, 3},
		{`last([1], [2]);`, "wrong number of arguments. got=2, want=1"},
		{`last("foo");`, "argument to 'last' must be an array, got STRING"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		switch expected := tt.expected.(type) {
		case int:
			testIntegerObject(tt.input, t, evaluated, int64(expected))
		case []int64:
			testIntegerArray(tt.input, t, evaluated, expected)
		case string:
			errObj, ok := evaluated.(*object.Error)
			if !ok {
				t.Errorf("object is not Error. got=%T(%+v)", expected, expected)
				continue
			}
			if errObj.Message != expected {
				t.Errorf("wrong error message, expected=%q, got=%q", expected, errObj.Message)
			}
		}
	}
}

func TestArrayIndexExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"[1, 2, 3][0]", 1},
		{"[1, 2, 3][1]", 2},
		{"[1, 2, 3][2]", 3},
		{"let i = 0; [1][i];", 1},
		{"[1, 2, 3][1 + 1];", 3},
		{"let myArray = [1, 2, 3]; myArray[2];", 3},
		{"let myArray = [1, 2, 3]; myArray[0] + myArray[1] + myArray[2];", 6},
		{"let myArray = [1, 2, 3]; let i = myArray[0]; myArray[i]", 2},
		{"[1,2,3][3]", NULL},
		{"[1,2,3][-1]", NULL},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(tt.input, t, evaluated, int64(integer))
		} else {
			testNullObject(tt.input, t, evaluated)
		}
	}

}

func TestHashLiterals(t *testing.T) {
	input := `let two = "two";
{"one": 10 - 9, two: 1 + 1, "thr" + "ee": 6/2, 4: 4, true: 5, false: 6 }`

	evaluated := testEval(input)
	result, ok := evaluated.(*object.Hash)
	if !ok {
		t.Fatalf("eval didn't return Hash. got=%T(%+v)", evaluated, evaluated)
	}

	expected := map[object.HashKey]int64{
		(&object.String{Value: "one"}).HashKey():   1,
		(&object.String{Value: "two"}).HashKey():   2,
		(&object.String{Value: "three"}).HashKey(): 3,
		(&object.Integer{Value: 4}).HashKey():      4,
		TRUE.HashKey():                             5,
		FALSE.HashKey():                            6,
	}

	if len(result.Pairs) != len(expected) {
		t.Fatalf("Hash has wrong number of pairs. got=%d", len(result.Pairs))
	}

	for expectedKey, expectedValue := range expected {
		pairs, ok := result.Pairs[expectedKey]
		if !ok {
			t.Errorf("no pair for given key in Pairs")
		}
		testIntegerObject(pairs.Key.Inspect(), t, pairs.Value, expectedValue)
	}
}

func TestHashIndexExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`{"foo": 5}["foo"]`, 5},
		{`{"foo": 5}["bar"]`, nil},
		{`let key = "foo"; {"foo": 5}[key]`, 5},
		{`{}["foo"]`, nil},
		{`{true: 5}[true]`, 5},
		{`{false: 5}[false]`, 5},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(tt.input, t, evaluated, int64(integer))
		} else {
			testNullObject(tt.input, t, evaluated)
		}
	}
}
