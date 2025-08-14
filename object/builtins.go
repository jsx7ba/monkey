package object

import (
	"cmp"
	"context"
	"fmt"
	"os/exec"
	"strings"
)

var (
	NULL  = &Null{}
	TRUE  = &Boolean{Value: true}
	FALSE = &Boolean{Value: false}
)

var Builtins = map[string]*Builtin{
	"len":   {Fn: length},
	"puts":  {Fn: puts},
	"first": {Fn: first},
	"last":  {Fn: last},
	"push":  {Fn: push},
	"rest":  {Fn: rest},
	"exec":  {Fn: execFn},
	"cmp":   {Fn: cmpFn},
}

func length(args ...Object) Object {
	if len(args) != 1 {
		return newError("wrong number of arguments, got=%d, want=1", len(args))
	}
	switch arg := args[0].(type) {
	case *Array:
		return &Integer{Value: int64(len(arg.Elements))}
	case *String:
		return &Integer{Value: int64(len(arg.Value))}
	case *Hash:
		return &Integer{Value: int64(len(arg.Pairs))}
	default:
		return newError("argument to 'len' not supported, got %s", args[0].Type())
	}
}

func puts(args ...Object) Object {
	for _, v := range args {
		fmt.Println(v.Inspect())
	}
	return NULL
}

func first(args ...Object) Object {
	if err, ok := checkArray("first", args); !ok {
		return err
	}
	arr := args[0].(*Array)
	if len(arr.Elements) > 0 {
		return arr.Elements[0]
	}
	return NULL
}

func last(args ...Object) Object {
	if err, ok := checkArray("last", args); !ok {
		return err
	}
	arr := args[0].(*Array)
	if len(arr.Elements) > 0 {
		return arr.Elements[len(arr.Elements)-1]
	}
	return NULL
}

func rest(args ...Object) Object {
	if err, ok := checkArray("rest", args); !ok {
		return err
	}
	arr := args[0].(*Array)
	count := len(arr.Elements)
	if count > 0 {
		elems := make([]Object, count-1)
		copy(elems, arr.Elements[1:])
		return &Array{Elements: elems}
	}
	return &Array{Elements: make([]Object, 0)}
}

func push(args ...Object) Object {
	if len(args) != 2 {
		return newError("wrong number of arguments. got=%d, want=2", len(args))
	}
	if args[0].Type() != ARRAY_OBJ {
		return newError("argument to `push` must be ARRAY, got %s", args[0].Type())
	}

	arr := args[0].(*Array)
	newElements := append(arr.Elements, args[1])
	return &Array{Elements: newElements}
}

func checkArray(name string, args []Object) (Object, bool) {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1", len(args)), false
	}
	if args[0].Type() != ARRAY_OBJ {
		return newError("argument to '%s' must be an array, got %s", name, args[0].Type()), false
	}

	return nil, true
}

func execFn(args ...Object) Object {
	if args[0].Type() != STRING_OBJ {
		return newError("exec requires a string argument")
	}
	strObj := args[0].(*String)
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

	return &String{Value: string(out)}
}

func cmpFn(args ...Object) Object {
	if len(args) != 2 {
		return newError("cmp requires 2 arguments")
	}

	switch args[0].Type() {
	case FLOAT_OBJ:
		a, ok := args[0].(*Float)
		if !ok {
			return newError("first argument must be a float")
		}
		b, ok := args[1].(*Float)

		if !ok {
			return newError("second argument must be a float")
		}

		return &Integer{Value: int64(cmp.Compare(a.Value, b.Value))}
	case STRING_OBJ:
		a, ok := args[0].(*String)
		if !ok {
			return newError("first argument must be a string")
		}
		b, ok := args[1].(*String)
		if !ok {
			return newError("second argument must be a string")
		}
		return &Integer{Value: int64(strings.Compare(a.Value, b.Value))}
	default:
		return newError("unsupported type: %s", args[0].Type())
	}
}

func newError(format string, a ...interface{}) *Error {
	return &Error{Message: fmt.Sprintf(format, a...)}
}

func GetBuiltinByName(name string) *Builtin {
	fn, ok := Builtins[name]
	if ok {
		return fn
	}
	return nil
}
