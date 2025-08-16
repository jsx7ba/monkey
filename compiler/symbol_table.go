package compiler

type SymbolScope string

const (
	GlobalScope  SymbolScope = "GLOBAL"
	LocalScope   SymbolScope = "LOCAL"
	BuiltinScope SymbolScope = "BUILTIN"
)

type Symbol struct {
	Name  string
	Scope SymbolScope
	Index int
}

type SymbolTable struct {
	outer          *SymbolTable
	store          map[string]Symbol
	numDefinitions int
}

func NewSymbolTable() *SymbolTable {
	s := make(map[string]Symbol)
	return &SymbolTable{store: s}
}

func NewEnclosedSymbolTable(outer *SymbolTable) *SymbolTable {
	s := NewSymbolTable()
	s.outer = outer
	return s
}

func (s *SymbolTable) Define(name string) Symbol {
	scope := LocalScope
	if s.outer == nil {
		scope = GlobalScope
	}
	symbol := Symbol{Name: name, Index: s.numDefinitions, Scope: scope}
	s.store[name] = symbol
	s.numDefinitions++
	return symbol
}

func (s *SymbolTable) Resolve(name string) (Symbol, bool) {
	sym, ok := s.store[name]
	if ok {
		return sym, ok
	}

	currentStore := s.outer
	for !ok && currentStore != nil {
		sym, ok = currentStore.store[name]
		currentStore = currentStore.outer
	}

	return sym, ok
}

func (s *SymbolTable) DefineBuiltin(index int, name string) Symbol {
	symbol := Symbol{name, BuiltinScope, index}
	s.store[name] = symbol
	return symbol
}
