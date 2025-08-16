package evaluator

import "monkey/object"

var builtins = map[string]*object.Builtin{
	"len":   object.GetBuiltinByName("len"),
	"puts":  object.GetBuiltinByName("puts"),
	"first": object.GetBuiltinByName("first"),
	"last":  object.GetBuiltinByName("last"),
	"push":  object.GetBuiltinByName("push"),
	"rest":  object.GetBuiltinByName("rest"),
	"exec":  object.GetBuiltinByName("exec"),
	"cmp":   object.GetBuiltinByName("cmp"),
}
