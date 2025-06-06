package object

import (
	"bytes"
	"fmt"
	"hash/fnv"
	"monkey/ast"
	"strings"
)

type ObjectType int

const (
	INTEGER_OBJ ObjectType = iota + 1
	FLOAT_OBJ
	BOOLEAN_OBJ
	NULL_OBJ
	RETURN_VALUE_OBJ
	ERROR_OBJ
	FUNCTION_OBJ
	STRING_OBJ
	BUILTIN_OBJ
	ARRAY_OBJ
	HASH_OBJ
)

func (o ObjectType) String() string {
	name := ""
	switch o {
	case INTEGER_OBJ:
		name = "INTEGER"
	case FLOAT_OBJ:
		name = "FLOAT"
	case BOOLEAN_OBJ:
		name = "BOOLEAN"
	case NULL_OBJ:
		name = "NULL"
	case RETURN_VALUE_OBJ:
		name = "RETURN_VALUE"
	case ERROR_OBJ:
		name = "ERROR"
	case FUNCTION_OBJ:
		name = "FUNCTION"
	case STRING_OBJ:
		name = "STRING"
	case BUILTIN_OBJ:
		name = "BUILTIN"
	case ARRAY_OBJ:
		name = "ARRAY"
	case HASH_OBJ:
		name = "HASH"
	default:
		name = "unknown object type"
	}
	return name
}

type Hashable interface {
	HashKey() HashKey
}

type HashKey struct {
	Type  ObjectType
	Value uint64
}

type HashPair struct {
	Key   Object
	Value Object
}

type Hash struct {
	Pairs map[HashKey]HashPair
}

func (h *Hash) Type() ObjectType { return HASH_OBJ }
func (h *Hash) Inspect() string {
	var out bytes.Buffer
	pairs := []string{}

	for _, pair := range h.Pairs {
		pairs = append(pairs, fmt.Sprintf("%s: %s", pair.Key.Inspect(), pair.Value.Inspect()))
	}

	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")
	return out.String()
}

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Float struct {
	Value float64
}

func (f *Float) Inspect() string  { return fmt.Sprintf("%f", f.Value) }
func (f *Float) Type() ObjectType { return FLOAT_OBJ }
func (f *Float) HashKey() HashKey {
	return HashKey{Type: f.Type(), Value: uint64(f.Value)}
}

type Integer struct {
	Value int64
}

func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }
func (i *Integer) Type() ObjectType { return INTEGER_OBJ }
func (i *Integer) HashKey() HashKey {
	return HashKey{Type: i.Type(), Value: uint64(i.Value)}
}

type Boolean struct {
	Value bool
}

func (b *Boolean) Inspect() string  { return fmt.Sprintf("%t", b.Value) }
func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }
func (b *Boolean) HashKey() HashKey {
	var value uint64
	if b.Value {
		value = 1
	} else {
		value = 0
	}
	return HashKey{Type: b.Type(), Value: value}
}

type Null struct{}

func (n *Null) Inspect() string  { return "null" }
func (n *Null) Type() ObjectType { return NULL_OBJ }

type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Type() ObjectType { return RETURN_VALUE_OBJ }
func (rv *ReturnValue) Inspect() string  { return rv.Value.Inspect() }

type Error struct {
	Message string
}

func (e *Error) Type() ObjectType { return ERROR_OBJ }
func (e *Error) Inspect() string  { return "ERROR: " + e.Message }

type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

func (f *Function) Type() ObjectType { return FUNCTION_OBJ }
func (f *Function) Inspect() string {
	var out bytes.Buffer
	params := []string{}
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}
	out.WriteString("fn")
	out.WriteString("{")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString("} \n")
	out.WriteString(f.Body.String())
	out.WriteString("\n")
	return out.String()
}

func NewStringHashKey(value string) HashKey {
	h := fnv.New64a()
	_, _ = h.Write([]byte(value))
	return HashKey{Type: STRING_OBJ, Value: h.Sum64()}
}

type String struct {
	Value string
}

func (s *String) Type() ObjectType { return STRING_OBJ }
func (s *String) Inspect() string  { return s.Value }
func (s *String) HashKey() HashKey {
	return NewStringHashKey(s.Value)
}

type BuiltinFunction func(args ...Object) Object
type Builtin struct {
	Fn   BuiltinFunction
	Void bool
}

func (b *Builtin) Type() ObjectType { return BUILTIN_OBJ }
func (b *Builtin) Inspect() string  { return "builtin function" }

type Array struct {
	Elements []Object
}

func (ao *Array) Type() ObjectType { return ARRAY_OBJ }
func (ao *Array) Inspect() string {
	var out bytes.Buffer

	elements := []string{}
	for _, e := range ao.Elements {
		elements = append(elements, e.Inspect())
	}
	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")
	return out.String()
}
