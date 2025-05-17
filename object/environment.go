package object

import (
	"os"
	"strings"
)

var env map[HashKey]HashPair

func init() {
	env = make(map[HashKey]HashPair)
	for _, v := range os.Environ() {
		key, value, _ := strings.Cut(v, "=")
		env[NewStringHashKey(key)] =
			HashPair{
				Key:   &String{Value: key},
				Value: &String{Value: value},
			}
	}
}

type Environment struct {
	store map[string]Object
	outer *Environment
}

func NewEnvironment() *Environment {
	s := make(map[string]Object)
	// add default variables to the environment
	s["ENV"] = &Hash{Pairs: env}
	return &Environment{store: s}
}

func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

func (e *Environment) Get(name string) (Object, bool) {
	val, ok := e.store[name]
	if !ok && e.outer != nil {
		val, ok = e.outer.Get(name)
	}
	return val, ok
}

func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}
