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

		for i, constant := range comp.Bytecode().Constants {
			fmt.Printf("CONSTANT %d %p (%T): \n", i, constant, constant)
			switch c := constant.(type) {
			case *object.CompiledFunction:
				fmt.Printf("Instructions:\n%s", c.Instructions)
			case *object.Integer:
				fmt.Printf("Value:\n%d", c.Value)
			}
			fmt.Println()
		}

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
	case *object.Error:
		errObj, ok := actual.(*object.Error)
		if !ok {
			t.Errorf("[%s] object is not Error: %T(%+v)", input, actual, actual)
			return
		}
		if errObj.Message != e.Message {
			t.Errorf("[%s] wrong error message, expected %q, got =%q", input, e.Message, errObj.Message)
		}
	case map[object.HashKey]int64:
		hash, ok := actual.(*object.Hash)
		if !ok {
			t.Errorf("[%s] object is not Hash.  got=%T (%+v)", input, actual, actual)
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

func TestIndexExpressions(t *testing.T) {
	tests := []vmTestCase{
		{"[1,2,3][1]", 2},
		{"[1,2,3][0+2]", 3},
		{"[][0]", Null},
		{"[1,2,3][99]", Null},
		{"[1][-1]", Null},
		{"{1:1,2:2}[1]", 1},
		{"{1:1,2:2}[2]", 2},
		{"{1:1}[0]", Null},
		{"{}[0]", Null},
	}
	runVmTests(t, tests)
}

func TestCallingFunctionsWithoutArguments(t *testing.T) {
	tests := []vmTestCase{
		{
			`let fivePlusTen = fn() { 5 + 10 }; fivePlusTen()`, 15,
		},
		{
			`let one = fn() { 1; }; let two = fn(){ 2;}; one() + two();`, 3,
		},
		{
			`let a = fn() { 1; }; let b = fn(){ a() + 1;}; let c = fn() { b() + 1 }; c()`, 3,
		},
	}

	runVmTests(t, tests)
}

func TestFunctionsWithReturnStatement(t *testing.T) {
	tests := []vmTestCase{
		{`let earlyExit = fn() { return 99; 100; }; earlyExit()`, 99},
		{`let earlyExit = fn() { return 99; return 100; }; earlyExit()`, 99},
	}
	runVmTests(t, tests)
}

func TestFunctionsWithoutReturnValue(t *testing.T) {
	tests := []vmTestCase{
		{
			`let noReturn = fn() {}`, Null,
		},
	}
	runVmTests(t, tests)
}

func TestFirstClassFunctions(t *testing.T) {
	tests := []vmTestCase{
		{input: `let returnsOne = fn(){ 1; };
                 let returnsOneReturner = fn() { returnsOne; };
                 returnsOneReturner()();`,
			expected: 1,
		},
	}
	runVmTests(t, tests)
}

func TestCallingFunctionsWithBindings(t *testing.T) {
	tests := []vmTestCase{
		{
			input:    `let one = fn() { let one = 1; one }; one();`,
			expected: 1,
		},
		{
			input:    `let oneAndTwo = fn() { let two = 2; let one = 1; two + one; }; oneAndTwo();`,
			expected: 3,
		},
		{
			input: `let oneAndTwo = fn() { let one = 1; let two = 2; one + two; };
		          let threeAndFour = fn() { let three = 3; let four = 4; three + four; };
		           oneAndTwo() + threeAndFour();`,
			expected: 10,
		},
		{
			input: `let firstFoobar = fn() { let foobar = 50; foobar; };
		          let secondFoobar = fn() { let foobar = 100; foobar; };
		          firstFoobar() + secondFoobar();`,
			expected: 150,
		},
		{
			input: `let globalSeed = 50;
				    let minusOne = fn() {let num = 1; globalSeed - num;};
					let minusTwo = fn() { let num = 2; globalSeed - num;};
					minusOne() + minusTwo()`,
			expected: 97,
		},
	}
	runVmTests(t, tests)
}

func TestCallingFunctionsWithArgumentsAndBindings(t *testing.T) {
	tests := []vmTestCase{
		{
			input:    `let identity = fn(a) { a; }; identity(4);`,
			expected: 4,
		},
		{
			input:    `let sum = fn(a, b) { a + b; }; sum(1,2);`,
			expected: 3,
		},
		{
			input:    `let sum = fn(a, b) { let c = a + b; c; }; sum(1, 2) + sum(3, 4);`,
			expected: 10,
		},
		{
			input:    `let sum = fn(a, b) { let c = a + b; c; }; let outer = fn() { sum(1, 2) + sum(3, 4); }; outer();`,
			expected: 10,
		},
		{
			input: `let globalNum = 10; let sum = fn(a, b) { let c = a + b; c + globalNum;}; 
					let outer = fn() {  sum(1, 2) + sum(3, 4) + globalNum; };
					outer() + globalNum;`,
			expected: 50,
		},
	}
	runVmTests(t, tests)
}

func TestCallingFunctionsWithWrongArguments(t *testing.T) {
	tests := []vmTestCase{
		{
			input:    `fn() { 1; }(1);`,
			expected: `wrong number of arguments: want=0, got=1`,
		},
		{
			input:    `fn(a) { a; }();`,
			expected: `wrong number of arguments: want=1, got=0`,
		},
		{
			input:    `fn(a, b) { a, b; }(1);`,
			expected: `wrong number of arguments: want=2, got=1`,
		},
	}

	for _, tt := range tests {
		program := parse(tt.input)
		comp := compiler.New()
		err := comp.Compile(program)
		if err != nil {
			t.Fatalf("compiler error: %+v", err)
		}
		vm := New(comp.Bytecode())
		err = vm.Run()

		if err == nil {
			t.Fatalf("expected VM error but got success")
		}
		if err.Error() != tt.expected {
			t.Fatalf("wrong VM errorf: want=%q, got=%q", tt.expected, err)
		}
	}
}

func TestBuiltinFunctions(t *testing.T) {
	tests := []vmTestCase{
		{`len("")`, 0},
		{`len("four")`, 4},
		{`len("Hello World")`, 11},
		{`len(1)`, &object.Error{Message: "argument to 'len' not supported, got INTEGER"}},
		{`len("one", "two")`, &object.Error{Message: "wrong number of arguments, got=2, want=1"}},
		{`len([1,2,3])`, 3},
		{`len([])`, 0},
		{`first([1,2,3])`, 1},
		{`first([])`, Null},
		{`first(1)`, &object.Error{Message: "argument to 'first' must be ARRAY, got INTEGER"}},
		{`last([1,2,3])`, 3},
		{`last([])`, &Null},
		{`last(1)`, &object.Error{Message: "argument to 'last' must be ARRAY, got INTEGER"}},
		{`rest([1,2,3])`, []int{2, 3}},
		{`rest([])`, Null},
		{`push([], 1)`, []int{1}},
		{`push(1, 2)`, &object.Error{Message: "argument to 'push' must be ARRAY, got INTEGER"}},
	}

	runVmTests(t, tests)
}

func TestClosures(t *testing.T) {
	tests := []vmTestCase{
		{
			`let newClosure = fn(a) { fn() { a; } };
					let closure = newClosure(99);
					closure();`,
			99,
		},
		{
			`let newAdder = fn(a, b) { fn(c) { a + b + c } }; let adder = newAdder(1, 2); adder(8);`,
			11,
		},
		{
			`let newAdder = fn(a, b) { let c = a + b; fn(d) { c + d }; }; let adder = newAdder(1,2); adder(8);`,
			11,
		},
		{
			`let newAdderOuter = fn(a, b) { 
					let c = a + b;
					fn(d) { 
						let e = d + c;
						fn(f) { e + f; };
				    };
                  }; 
				let newAdderInner = newAdderOuter(1,2);
				let adder = newAdderInner(3);
				adder(8); `,
			14,
		},
		{
			`let a = 1;
				let newAdderOuter = fn(b) { 
					fn(c) {
						fn(d) { a + b + c + d };
		            };
				};
				let newAdderInner = newAdderOuter(2);
				let adder = newAdderInner(3);
				adder(8);`,
			14,
		},
		{
			`let newClosure = fn(a, b) {
						let one = fn() { a; };
					    let two = fn() { b; };
						fn() { one() + two(); };
					};
                   let closure = newClosure(9, 90);
                   closure();`,
			99,
		},
	}
	runVmTests(t, tests)
}

func TestRecursiveFunctions(t *testing.T) {
	tests := []vmTestCase{
		{
			`let countDown = fn(x) { 
					if (x == 0) { return 0; } else { countDown(x -1); } 
                   };
				countDown(1);`,
			0,
		},
		{
			`let countDown = fn(x) { if (x == 0) { return 0; } else { countDown(x -1); } };
					let wrapper = fn() { countDown(1); };
					wrapper();`,
			0,
		},
		{
			`let wrapper = fn() {
						let countDown = fn(x) { if (x == 0) { return 0; } else { countDown(x - 1); } }; 
					countDown(1); }; 
                    wrapper(); `,
			0,
		},
	}

	runVmTests(t, tests)
}
