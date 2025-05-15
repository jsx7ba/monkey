package evaluator

import (
	"fmt"
	"monkey/object"
)

var builtins = map[string]*object.Builtin{
	"len":   {Fn: monkeyLen},
	"first": {Fn: first},
	"last":  {Fn: last},
	"rest":  {Fn: rest},
	"push":  {Fn: push},
	"puts":  {Fn: puts, Void: true},
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
