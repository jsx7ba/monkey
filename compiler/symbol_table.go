package compiler

type SymbolScope string

const (
	GlobalScope   SymbolScope = "GLOBAL"
	LocalScope    SymbolScope = "LOCAL"
	BuiltinScope  SymbolScope = "BUILTIN"
	FreeScope     SymbolScope = "FREE"
	FunctionScope SymbolScope = "FUNCTION"
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
	FreeSymbols    []Symbol
}

func NewSymbolTable() *SymbolTable {
	s := make(map[string]Symbol)
	return &SymbolTable{store: s, FreeSymbols: make([]Symbol, 0)}
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
	if !ok && s.outer != nil {
		sym, ok = s.outer.Resolve(name)
		if !ok {
			return sym, ok
		}
		if sym.Scope == GlobalScope || sym.Scope == BuiltinScope {
			return sym, ok
		}
		free := s.defineFree(sym)
		return free, true
	}
	return sym, ok
}

func (s *SymbolTable) DefineBuiltin(index int, name string) Symbol {
	symbol := Symbol{name, BuiltinScope, index}
	s.store[name] = symbol
	return symbol
}

func (s *SymbolTable) defineFree(original Symbol) Symbol {
	s.FreeSymbols = append(s.FreeSymbols, original)
	symbol := Symbol{
		Name:  original.Name,
		Scope: FreeScope,
		Index: len(s.FreeSymbols) - 1,
	}

	s.store[original.Name] = symbol
	return symbol
}

func (s *SymbolTable) DefineFunctionName(name string) Symbol {
	symbol := Symbol{name, FunctionScope, 0}
	s.store[name] = symbol
	return symbol
}
