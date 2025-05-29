package lexer

import (
	"fmt"
	"monkey/token"
	"testing"
)

func LineInfoEquals(expected, actual token.LineInfo) bool {
	return expected.FileName == actual.FileName &&
		expected.Line == actual.Line &&
		expected.Char == expected.Char
}

func TestNextToken(t *testing.T) {
	input := `let five = 55;
let ten = 10;

let add = fn(x,y) {
	x + y;
};

let result = add(five, ten);
!-/*5;
5 < 10 > 5;

if (5 < 10) { 
	return true;
} else { 
	return false;
}

10 == 10;
10 != 9;
"foobar"
"foo bar"
[1, 2];
{"foo": "bar"}
let hexNum = 0x33;
let octNum = 033;
let fa = 3.141592;
let fb = 0.003;
let array = ["a", "b"];
let xx = array[0];
`

	tests := []struct {
		expectedType     token.TokenType
		expectedLiteral  string
		expectedLineInfo token.LineInfo
	}{
		{token.LET, "let", token.LineInfo{FileName: "REPL", Line: 1, Char: 1}},
		{token.IDENT, "five", token.LineInfo{FileName: "REPL", Line: 1, Char: 5}},
		{token.ASSIGN, "=", token.LineInfo{FileName: "REPL", Line: 1, Char: 10}},
		{token.INT, "55", token.LineInfo{FileName: "REPL", Line: 1, Char: 12}},
		{token.SEMICOLON, ";", token.LineInfo{FileName: "REPL", Line: 1, Char: 15}},

		{token.LET, "let", token.LineInfo{FileName: "REPL", Line: 2, Char: 1}},
		{token.IDENT, "ten", token.LineInfo{FileName: "REPL", Line: 2, Char: 5}},
		{token.ASSIGN, "=", token.LineInfo{FileName: "REPL", Line: 2, Char: 9}},
		{token.INT, "10", token.LineInfo{FileName: "REPL", Line: 2, Char: 11}},
		{token.SEMICOLON, ";", token.LineInfo{FileName: "REPL", Line: 2, Char: 13}},

		{token.LET, "let", token.LineInfo{FileName: "REPL", Line: 4, Char: 1}},
		{token.IDENT, "add", token.LineInfo{FileName: "REPL", Line: 4, Char: 5}},
		{token.ASSIGN, "=", token.LineInfo{FileName: "REPL", Line: 4, Char: 9}},
		{token.FUNCTION, "fn", token.LineInfo{FileName: "REPL", Line: 4, Char: 11}},
		{token.LPAREN, "(", token.LineInfo{FileName: "REPL", Line: 4, Char: 13}},
		{token.IDENT, "x", token.LineInfo{FileName: "REPL", Line: 4, Char: 14}},
		{token.COMMA, ",", token.LineInfo{FileName: "REPL", Line: 4, Char: 15}},
		{token.IDENT, "y", token.LineInfo{FileName: "REPL", Line: 4, Char: 16}},
		{token.RPAREN, ")", token.LineInfo{FileName: "REPL", Line: 4, Char: 17}},
		{token.LBRACE, "{", token.LineInfo{FileName: "REPL", Line: 4, Char: 19}},
		{token.IDENT, "x", token.LineInfo{FileName: "REPL", Line: 5, Char: 2}},
		{token.PLUS, "+", token.LineInfo{FileName: "REPL", Line: 5, Char: 4}},
		{token.IDENT, "y", token.LineInfo{FileName: "REPL", Line: 5, Char: 6}},
		{token.SEMICOLON, ";", token.LineInfo{FileName: "REPL", Line: 5, Char: 7}},
		{token.RBRACE, "}", token.LineInfo{FileName: "REPL", Line: 6, Char: 1}},
		{token.SEMICOLON, ";", token.LineInfo{FileName: "REPL", Line: 6, Char: 2}},

		{token.LET, "let", token.LineInfo{FileName: "REPL", Line: 8, Char: 1}},
		{token.IDENT, "result", token.LineInfo{FileName: "REPL", Line: 8, Char: 5}},
		{token.ASSIGN, "=", token.LineInfo{FileName: "REPL", Line: 8, Char: 12}},
		{token.IDENT, "add", token.LineInfo{FileName: "REPL", Line: 8, Char: 14}},
		{token.LPAREN, "(", token.LineInfo{FileName: "REPL", Line: 8, Char: 17}},
		{token.IDENT, "five", token.LineInfo{FileName: "REPL", Line: 8, Char: 18}},
		{token.COMMA, ",", token.LineInfo{FileName: "REPL", Line: 8, Char: 22}},
		{token.IDENT, "ten", token.LineInfo{FileName: "REPL", Line: 8, Char: 24}},
		{token.RPAREN, ")", token.LineInfo{FileName: "REPL", Line: 8, Char: 27}},
		{token.SEMICOLON, ";", token.LineInfo{FileName: "REPL", Line: 8, Char: 28}},

		{token.BANG, "!", token.LineInfo{FileName: "REPL", Line: 9, Char: 1}},
		{token.MINUS, "-", token.LineInfo{FileName: "REPL", Line: 9, Char: 2}},
		{token.SLASH, "/", token.LineInfo{FileName: "REPL", Line: 9, Char: 3}},
		{token.ASTERISK, "*", token.LineInfo{FileName: "REPL", Line: 9, Char: 4}},
		{token.INT, "5", token.LineInfo{FileName: "REPL", Line: 9, Char: 5}},
		{token.SEMICOLON, ";", token.LineInfo{FileName: "REPL", Line: 9, Char: 6}},

		{token.INT, "5", token.LineInfo{FileName: "REPL", Line: 10, Char: 1}},
		{token.LT, "<", token.LineInfo{FileName: "REPL", Line: 10, Char: 3}},
		{token.INT, "10", token.LineInfo{FileName: "REPL", Line: 10, Char: 5}},
		{token.GT, ">", token.LineInfo{FileName: "REPL", Line: 10, Char: 8}},
		{token.INT, "5", token.LineInfo{FileName: "REPL", Line: 10, Char: 10}},
		{token.SEMICOLON, ";", token.LineInfo{FileName: "REPL", Line: 10, Char: 11}},

		{token.IF, "if", token.LineInfo{FileName: "REPL", Line: 12, Char: 1}},
		{token.LPAREN, "(", token.LineInfo{FileName: "REPL", Line: 12, Char: 4}},
		{token.INT, "5", token.LineInfo{FileName: "REPL", Line: 12, Char: 5}},
		{token.LT, "<", token.LineInfo{FileName: "REPL", Line: 12, Char: 7}},
		{token.INT, "10", token.LineInfo{FileName: "REPL", Line: 12, Char: 9}},
		{token.RPAREN, ")", token.LineInfo{FileName: "REPL", Line: 12, Char: 11}},
		{token.LBRACE, "{", token.LineInfo{FileName: "REPL", Line: 12, Char: 13}},
		{token.RETURN, "return", token.LineInfo{FileName: "REPL", Line: 13, Char: 2}},
		{token.TRUE, "true", token.LineInfo{FileName: "REPL", Line: 13, Char: 9}},
		{token.SEMICOLON, ";", token.LineInfo{FileName: "REPL", Line: 13, Char: 13}},
		{token.RBRACE, "}", token.LineInfo{FileName: "REPL", Line: 14, Char: 1}},
		{token.ELSE, "else", token.LineInfo{FileName: "REPL", Line: 14, Char: 3}},
		{token.LBRACE, "{", token.LineInfo{FileName: "REPL", Line: 14, Char: 8}},
		{token.RETURN, "return", token.LineInfo{FileName: "REPL", Line: 15, Char: 2}},
		{token.FALSE, "false", token.LineInfo{FileName: "REPL", Line: 15, Char: 9}},
		{token.SEMICOLON, ";", token.LineInfo{FileName: "REPL", Line: 15, Char: 14}},
		{token.RBRACE, "}", token.LineInfo{FileName: "REPL", Line: 16, Char: 1}},

		{token.INT, "10", token.LineInfo{FileName: "REPL", Line: 18, Char: 1}},
		{token.EQ, "==", token.LineInfo{FileName: "REPL", Line: 18, Char: 4}},
		{token.INT, "10", token.LineInfo{FileName: "REPL", Line: 18, Char: 7}},
		{token.SEMICOLON, ";", token.LineInfo{FileName: "REPL", Line: 18, Char: 9}},

		{token.INT, "10", token.LineInfo{FileName: "REPL", Line: 19, Char: 1}},
		{token.NOT_EQ, "!=", token.LineInfo{FileName: "REPL", Line: 19, Char: 4}},
		{token.INT, "9", token.LineInfo{FileName: "REPL", Line: 19, Char: 7}},
		{token.SEMICOLON, ";", token.LineInfo{FileName: "REPL", Line: 19, Char: 8}},
		{token.STRING, "foobar", token.LineInfo{FileName: "REPL", Line: 20, Char: 1}},
		{token.STRING, "foo bar", token.LineInfo{FileName: "REPL", Line: 21, Char: 1}},

		{token.LBRACKET, "[", token.LineInfo{FileName: "REPL", Line: 22, Char: 2}},
		{token.INT, "1", token.LineInfo{FileName: "REPL", Line: 22, Char: 3}},
		{token.COMMA, ",", token.LineInfo{FileName: "REPL", Line: 22, Char: 5}},
		{token.INT, "2", token.LineInfo{FileName: "REPL", Line: 22, Char: 6}},
		{token.RBRACKET, "]", token.LineInfo{FileName: "REPL", Line: 22, Char: 7}},
		{token.SEMICOLON, ";", token.LineInfo{FileName: "REPL", Line: 22, Char: 8}},

		{token.LBRACE, "{", token.LineInfo{FileName: "REPL", Line: 23, Char: 2}},
		{token.STRING, "foo", token.LineInfo{FileName: "REPL", Line: 23, Char: 7}},
		{token.COLON, ":", token.LineInfo{FileName: "REPL", Line: 23, Char: 9}},
		{token.STRING, "bar", token.LineInfo{FileName: "REPL", Line: 23, Char: 14}},
		{token.RBRACE, "}", token.LineInfo{FileName: "REPL", Line: 23, Char: 14}},

		{token.LET, "let", token.LineInfo{FileName: "REPL", Line: 24, Char: 1}},
		{token.IDENT, "hexNum", token.LineInfo{FileName: "REPL", Line: 24, Char: 5}},
		{token.ASSIGN, "=", token.LineInfo{FileName: "REPL", Line: 24, Char: 12}},
		{token.INT, "0x33", token.LineInfo{FileName: "REPL", Line: 24, Char: 14}},
		{token.SEMICOLON, ";", token.LineInfo{FileName: "REPL", Line: 24, Char: 18}},

		{token.LET, "let", token.LineInfo{FileName: "REPL", Line: 25, Char: 1}},
		{token.IDENT, "octNum", token.LineInfo{FileName: "REPL", Line: 25, Char: 5}},
		{token.ASSIGN, "=", token.LineInfo{FileName: "REPL", Line: 25, Char: 12}},
		{token.INT, "033", token.LineInfo{FileName: "REPL", Line: 25, Char: 14}},
		{token.SEMICOLON, ";", token.LineInfo{FileName: "REPL", Line: 25, Char: 17}},

		{token.LET, "let", token.LineInfo{FileName: "REPL", Line: 26, Char: 1}},
		{token.IDENT, "fa", token.LineInfo{FileName: "REPL", Line: 26, Char: 5}},
		{token.ASSIGN, "=", token.LineInfo{FileName: "REPL", Line: 26, Char: 8}},
		{token.FLOAT, "3.141592", token.LineInfo{FileName: "REPL", Line: 26, Char: 10}},
		{token.SEMICOLON, ";", token.LineInfo{FileName: "REPL", Line: 26, Char: 18}},

		{token.LET, "let", token.LineInfo{FileName: "REPL", Line: 27, Char: 1}},
		{token.IDENT, "fb", token.LineInfo{FileName: "REPL", Line: 27, Char: 5}},
		{token.ASSIGN, "=", token.LineInfo{FileName: "REPL", Line: 27, Char: 8}},
		{token.FLOAT, "0.003", token.LineInfo{FileName: "REPL", Line: 27, Char: 10}},
		{token.SEMICOLON, ";", token.LineInfo{FileName: "REPL", Line: 27, Char: 15}},

		{token.LET, "let", token.LineInfo{FileName: "REPL", Line: 28, Char: 1}},
		{token.IDENT, "array", token.LineInfo{FileName: "REPL", Line: 28, Char: 5}},
		{token.ASSIGN, "=", token.LineInfo{FileName: "REPL", Line: 28, Char: 11}},
		{token.LBRACKET, "[", token.LineInfo{FileName: "REPL", Line: 28, Char: 13}},
		{token.STRING, "a", token.LineInfo{FileName: "REPL", Line: 28, Char: 14}},
		{token.COMMA, ",", token.LineInfo{FileName: "REPL", Line: 28, Char: 17}},
		{token.STRING, "b", token.LineInfo{FileName: "REPL", Line: 28, Char: 19}},
		{token.RBRACKET, "]", token.LineInfo{FileName: "REPL", Line: 28, Char: 22}},
		{token.SEMICOLON, ";", token.LineInfo{FileName: "REPL", Line: 28, Char: 23}},

		{token.LET, "let", token.LineInfo{FileName: "REPL", Line: 29, Char: 1}},
		{token.IDENT, "xx", token.LineInfo{FileName: "REPL", Line: 29, Char: 5}},
		{token.ASSIGN, "=", token.LineInfo{FileName: "REPL", Line: 29, Char: 8}},
		{token.IDENT, "array", token.LineInfo{FileName: "REPL", Line: 29, Char: 10}},
		{token.LBRACKET, "[", token.LineInfo{FileName: "REPL", Line: 29, Char: 15}},
		{token.INT, "0", token.LineInfo{FileName: "REPL", Line: 29, Char: 16}},
		{token.RBRACKET, "]", token.LineInfo{FileName: "REPL", Line: 29, Char: 17}},
		{token.SEMICOLON, ";", token.LineInfo{FileName: "REPL", Line: 29, Char: 18}},

		{token.EOF, "", token.LineInfo{FileName: "REPL", Line: 30, Char: 0}},
	}

	lex := NewFromString("REPL", input)

	for _, tt := range tests {
		tok := lex.NextToken()
		fmt.Printf("%+v\n", tok)
		if tok.Type != tt.expectedType {
			t.Fatalf("test[%s] - expected type %s, got %s", tt.expectedLiteral, tt.expectedType, tok.Type)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("test[%s] - expected literal %s, got %s", tt.expectedLiteral, tt.expectedLiteral, tok.Literal)
		}
		if !LineInfoEquals(tt.expectedLineInfo, tok.LineInfo) {
			t.Fatalf("test[%s] - expected [%s] got [%s]", tt.expectedLiteral, tt.expectedLineInfo, tok.LineInfo)
		}
	}
}

func TestLineNumbersWithComments(t *testing.T) {
	input := `0; # tail comment
   # line commented out
42;`

	l := NewFromString("REPL", input)
	l.NextToken() // read the 0
	l.NextToken() // read the ;
	tok := l.NextToken()
	if tok.Type != token.INT {
		t.Errorf("expected %s got %s", token.INT, tok.Type)
	}
	if tok.Literal != "42" {
		t.Errorf("expected %s got %s", "42", tok.Literal)
	}

	expectedLineInfo := token.LineInfo{
		FileName: "REPL",
		Line:     3,
		Char:     0,
	}
	if !LineInfoEquals(expectedLineInfo, tok.LineInfo) {
		t.Errorf("expected %s got %s", expectedLineInfo, tok.LineInfo)
	}
}
