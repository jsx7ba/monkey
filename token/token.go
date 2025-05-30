package token

import "fmt"

type LineInfo struct {
	FileName string
	Line     int
	Char     int
}

func (li LineInfo) String() string {
	return fmt.Sprintf("%s: line %d, char %d", li.FileName, li.Line, li.Char)
}

type TokenType string

type Token struct {
	Type     TokenType
	Literal  string
	LineInfo LineInfo
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// Identifiers
	IDENT = "IDENT"
	INT   = "INT"
	FLOAT = "FLOAT"

	// Operators
	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	BANG     = "!"
	ASTERISK = "*"
	SLASH    = "/"
	LT       = "<"
	GT       = ">"
	EQ       = "=="
	NOT_EQ   = "!="

	// Delimiters
	COMMA     = ","
	SEMICOLON = ";"
	LPAREN    = "("
	RPAREN    = ")"
	LBRACE    = "{"
	RBRACE    = "}"

	// Keywords
	FUNCTION = "FUNCTION"
	LET      = "LET"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	IF       = "IF"
	ELSE     = "ELSE"
	RETURN   = "RETURN"

	// Built-ins
	STRING   = "STRING"
	LBRACKET = "["
	RBRACKET = "]"
	COLON    = ":"

	// Comment
	COMMENT = "#"
)

var keywords = map[string]TokenType{
	"fn":     FUNCTION,
	"let":    LET,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
