package evaluator

import (
	"cmp"
	"context"
	"fmt"
	"monkey/object"
	"os/exec"
	"strings"
)

var builtins = map[string]*object.Builtin{
	"len":   {Fn: monkeyLen},
	"first": {Fn: first},
	"last":  {Fn: last},
	"rest":  {Fn: rest},
	"push":  {Fn: push},
	"puts":  {Fn: puts, Void: true},
	"exec":  {Fn: execFn},
	"cmp":   {Fn: cmpFn},
}

func monkeyLen(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1", len(args))
	}

	switch arg := args[0].(type) {
	case *object.String:
		return &object.Integer{Value: int64(len(arg.Value))}
	case *object.Array:
		return &object.Integer{Value: int64(len(arg.Elements))}
	default:
		return newError("argument to 'len' not supported, got %s", args[0].Type())
	}
}

func checkArray(name string, args []object.Object) (object.Object, bool) {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1", len(args)), false
	}
	if args[0].Type() != object.ARRAY_OBJ {
		return newError("argument to '%s' must be an array, got %s", name, args[0].Type()), false
	}

	return nil, true
}

func first(args ...object.Object) object.Object {
	if err, ok := checkArray("first", args); !ok {
		return err
	}
	arr := args[0].(*object.Array)
	if len(arr.Elements) > 0 {
		return arr.Elements[0]
	}
	return NULL
}

func last(args ...object.Object) object.Object {
	if err, ok := checkArray("last", args); !ok {
		return err
	}
	arr := args[0].(*object.Array)
	if len(arr.Elements) > 0 {
		return arr.Elements[len(arr.Elements)-1]
	}
	return NULL
}

func rest(args ...object.Object) object.Object {
	if err, ok := checkArray("rest", args); !ok {
		return err
	}

	arr := args[0].(*object.Array)
	count := len(arr.Elements)
	if count > 0 {
		elems := make([]object.Object, count-1)
		copy(elems, arr.Elements[1:])
		return &object.Array{Elements: elems}
	}
	return &object.Array{Elements: make([]object.Object, 0)}
}

func push(args ...object.Object) object.Object {
	if len(args) != 2 {
		return newError("wrong number of arguments. got=%d, want=2", len(args))
	}
	if args[0].Type() != object.ARRAY_OBJ {
		return newError("argument to `push` must be ARRAY, got %s", args[0].Type())
	}

	arr := args[0].(*object.Array)
	newElements := append(arr.Elements, args[1])
	return &object.Array{Elements: newElements}
}

func puts(args ...object.Object) object.Object {
	for _, o := range args {
		fmt.Println(o.Inspect())
	}
	return NULL
}

func execFn(args ...object.Object) object.Object {
	if args[0].Type() != object.STRING_OBJ {
		return newError("exec requires a string argument")
	}
	strObj := args[0].(*object.String)
	parts := strings.Split(strObj.Value, " ")
	cmd := exec.CommandContext(context.Background(), parts[0], parts[1:]...)
	err := cmd.Err
	if err != nil {
		return newError("%+v", err)
	}

	out, err := cmd.Output()
	if err != nil {
		return newError("exec failed: %+v", err)
	}

	return &object.String{Value: string(out)}
}

func cmpFn(args ...object.Object) object.Object {
	if len(args) != 2 {
		return newError("cmp requires 2 arguments")
	}

	switch args[0].Type() {
	case object.FLOAT_OBJ:
		a, ok := args[0].(*object.Float)
		if !ok {
			return newError("first argument must be a float")
		}
		b, ok := args[1].(*object.Float)

		if !ok {
			return newError("second argument must be a float")
		}

		return &object.Integer{Value: int64(cmp.Compare(a.Value, b.Value))}
	case object.STRING_OBJ:
		a, ok := args[0].(*object.String)
		if !ok {
			return newError("first argument must be a string")
		}
		b, ok := args[1].(*object.String)
		if !ok {
			return newError("second argument must be a string")
		}
		return &object.Integer{Value: int64(strings.Compare(a.Value, b.Value))}
	default:
		return newError("unsupported type: %s", args[0].Type())
	}
}
