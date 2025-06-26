package lexer

import (
	"fmt"
	"monkey/token"
	"testing"
)

func LineInfoEquals(expected, actual token.LineInfo) bool {
	return expected.FileIndex == actual.FileIndex &&
		expected.Line == actual.Line &&
		expected.Char == actual.Char
}

func TestNextToken(t *testing.T) {
	token.ResetForTesting()
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
{"foo": "bar"};
let hexNum = 0x33;
let octNum = 033;
let fa = 3.141592;
let fb = 0.003;
let array = ["a", "b"];
let xx = array[0];
`

	srcHandle := token.SourceHandle(0)

	tests := []struct {
		expectedType     token.TokenType
		expectedLiteral  string
		expectedLineInfo token.LineInfo
	}{
		{token.LET, "let", srcHandle.LineInfo(1, 1)},
		{token.IDENT, "five", srcHandle.LineInfo(1, 5)},
		{token.ASSIGN, "=", srcHandle.LineInfo(1, 10)},
		{token.INT, "55", srcHandle.LineInfo(1, 12)},
		{token.SEMICOLON, ";", srcHandle.LineInfo(1, 14)},

		{token.LET, "let", srcHandle.LineInfo(2, 1)},
		{token.IDENT, "ten", srcHandle.LineInfo(2, 5)},
		{token.ASSIGN, "=", srcHandle.LineInfo(2, 9)},
		{token.INT, "10", srcHandle.LineInfo(2, 11)},
		{token.SEMICOLON, ";", srcHandle.LineInfo(2, 13)},

		{token.LET, "let", srcHandle.LineInfo(4, 1)},
		{token.IDENT, "add", srcHandle.LineInfo(4, 5)},
		{token.ASSIGN, "=", srcHandle.LineInfo(4, 9)},
		{token.FUNCTION, "fn", srcHandle.LineInfo(4, 11)},
		{token.LPAREN, "(", srcHandle.LineInfo(4, 13)},
		{token.IDENT, "x", srcHandle.LineInfo(4, 14)},
		{token.COMMA, ",", srcHandle.LineInfo(4, 15)},
		{token.IDENT, "y", srcHandle.LineInfo(4, 16)},
		{token.RPAREN, ")", srcHandle.LineInfo(4, 17)},
		{token.LBRACE, "{", srcHandle.LineInfo(4, 19)},
		{token.IDENT, "x", srcHandle.LineInfo(5, 2)},
		{token.PLUS, "+", srcHandle.LineInfo(5, 4)},
		{token.IDENT, "y", srcHandle.LineInfo(5, 6)},
		{token.SEMICOLON, ";", srcHandle.LineInfo(5, 7)},
		{token.RBRACE, "}", srcHandle.LineInfo(6, 1)},
		{token.SEMICOLON, ";", srcHandle.LineInfo(6, 2)},

		{token.LET, "let", srcHandle.LineInfo(8, 1)},
		{token.IDENT, "result", srcHandle.LineInfo(8, 5)},
		{token.ASSIGN, "=", srcHandle.LineInfo(8, 12)},
		{token.IDENT, "add", srcHandle.LineInfo(8, 14)},
		{token.LPAREN, "(", srcHandle.LineInfo(8, 17)},
		{token.IDENT, "five", srcHandle.LineInfo(8, 18)},
		{token.COMMA, ",", srcHandle.LineInfo(8, 22)},
		{token.IDENT, "ten", srcHandle.LineInfo(8, 24)},
		{token.RPAREN, ")", srcHandle.LineInfo(8, 27)},
		{token.SEMICOLON, ";", srcHandle.LineInfo(8, 28)},

		{token.BANG, "!", srcHandle.LineInfo(9, 1)},
		{token.MINUS, "-", srcHandle.LineInfo(9, 2)},
		{token.SLASH, "/", srcHandle.LineInfo(9, 3)},
		{token.ASTERISK, "*", srcHandle.LineInfo(9, 4)},
		{token.INT, "5", srcHandle.LineInfo(9, 5)},
		{token.SEMICOLON, ";", srcHandle.LineInfo(9, 6)},

		{token.INT, "5", srcHandle.LineInfo(10, 1)},
		{token.LT, "<", srcHandle.LineInfo(10, 3)},
		{token.INT, "10", srcHandle.LineInfo(10, 5)},
		{token.GT, ">", srcHandle.LineInfo(10, 8)},
		{token.INT, "5", srcHandle.LineInfo(10, 10)},
		{token.SEMICOLON, ";", srcHandle.LineInfo(10, 11)},

		{token.IF, "if", srcHandle.LineInfo(12, 1)},
		{token.LPAREN, "(", srcHandle.LineInfo(12, 4)},
		{token.INT, "5", srcHandle.LineInfo(12, 5)},
		{token.LT, "<", srcHandle.LineInfo(12, 7)},
		{token.INT, "10", srcHandle.LineInfo(12, 9)},
		{token.RPAREN, ")", srcHandle.LineInfo(12, 11)},
		{token.LBRACE, "{", srcHandle.LineInfo(12, 13)},
		{token.RETURN, "return", srcHandle.LineInfo(13, 2)},
		{token.TRUE, "true", srcHandle.LineInfo(13, 9)},
		{token.SEMICOLON, ";", srcHandle.LineInfo(13, 13)},
		{token.RBRACE, "}", srcHandle.LineInfo(14, 1)},
		{token.ELSE, "else", srcHandle.LineInfo(14, 3)},
		{token.LBRACE, "{", srcHandle.LineInfo(14, 8)},
		{token.RETURN, "return", srcHandle.LineInfo(15, 2)},
		{token.FALSE, "false", srcHandle.LineInfo(15, 9)},
		{token.SEMICOLON, ";", srcHandle.LineInfo(15, 14)},
		{token.RBRACE, "}", srcHandle.LineInfo(16, 1)},

		{token.INT, "10", srcHandle.LineInfo(18, 1)},
		{token.EQ, "==", srcHandle.LineInfo(18, 4)},
		{token.INT, "10", srcHandle.LineInfo(18, 7)},
		{token.SEMICOLON, ";", srcHandle.LineInfo(18, 9)},

		{token.INT, "10", srcHandle.LineInfo(19, 1)},
		{token.NOT_EQ, "!=", srcHandle.LineInfo(19, 4)},
		{token.INT, "9", srcHandle.LineInfo(19, 7)},
		{token.SEMICOLON, ";", srcHandle.LineInfo(19, 8)},
		{token.STRING, "foobar", srcHandle.LineInfo(20, 1)},
		{token.STRING, "foo bar", srcHandle.LineInfo(21, 1)},

		{token.LBRACKET, "[", srcHandle.LineInfo(22, 1)},
		{token.INT, "1", srcHandle.LineInfo(22, 2)},
		{token.COMMA, ",", srcHandle.LineInfo(22, 3)},
		{token.INT, "2", srcHandle.LineInfo(22, 5)},
		{token.RBRACKET, "]", srcHandle.LineInfo(22, 6)},
		{token.SEMICOLON, ";", srcHandle.LineInfo(22, 7)},

		{token.LBRACE, "{", srcHandle.LineInfo(23, 1)},
		{token.STRING, "foo", srcHandle.LineInfo(23, 2)},
		{token.COLON, ":", srcHandle.LineInfo(23, 7)},
		{token.STRING, "bar", srcHandle.LineInfo(23, 9)},
		{token.RBRACE, "}", srcHandle.LineInfo(23, 14)},
		{token.SEMICOLON, ";", srcHandle.LineInfo(23, 15)},

		{token.LET, "let", srcHandle.LineInfo(24, 1)},
		{token.IDENT, "hexNum", srcHandle.LineInfo(24, 5)},
		{token.ASSIGN, "=", srcHandle.LineInfo(24, 12)},
		{token.INT, "0x33", srcHandle.LineInfo(24, 14)},
		{token.SEMICOLON, ";", srcHandle.LineInfo(24, 18)},

		{token.LET, "let", srcHandle.LineInfo(25, 1)},
		{token.IDENT, "octNum", srcHandle.LineInfo(25, 5)},
		{token.ASSIGN, "=", srcHandle.LineInfo(25, 12)},
		{token.INT, "033", srcHandle.LineInfo(25, 14)},
		{token.SEMICOLON, ";", srcHandle.LineInfo(25, 17)},

		{token.LET, "let", srcHandle.LineInfo(26, 1)},
		{token.IDENT, "fa", srcHandle.LineInfo(26, 5)},
		{token.ASSIGN, "=", srcHandle.LineInfo(26, 8)},
		{token.FLOAT, "3.141592", srcHandle.LineInfo(26, 10)},
		{token.SEMICOLON, ";", srcHandle.LineInfo(26, 18)},

		{token.LET, "let", srcHandle.LineInfo(27, 1)},
		{token.IDENT, "fb", srcHandle.LineInfo(27, 5)},
		{token.ASSIGN, "=", srcHandle.LineInfo(27, 8)},
		{token.FLOAT, "0.003", srcHandle.LineInfo(27, 10)},
		{token.SEMICOLON, ";", srcHandle.LineInfo(27, 15)},

		{token.LET, "let", srcHandle.LineInfo(28, 1)},
		{token.IDENT, "array", srcHandle.LineInfo(28, 5)},
		{token.ASSIGN, "=", srcHandle.LineInfo(28, 11)},
		{token.LBRACKET, "[", srcHandle.LineInfo(28, 13)},
		{token.STRING, "a", srcHandle.LineInfo(28, 14)},
		{token.COMMA, ",", srcHandle.LineInfo(28, 17)},
		{token.STRING, "b", srcHandle.LineInfo(28, 19)},
		{token.RBRACKET, "]", srcHandle.LineInfo(28, 22)},
		{token.SEMICOLON, ";", srcHandle.LineInfo(28, 23)},

		{token.LET, "let", srcHandle.LineInfo(29, 1)},
		{token.IDENT, "xx", srcHandle.LineInfo(29, 5)},
		{token.ASSIGN, "=", srcHandle.LineInfo(29, 8)},
		{token.IDENT, "array", srcHandle.LineInfo(29, 10)},
		{token.LBRACKET, "[", srcHandle.LineInfo(29, 15)},
		{token.INT, "0", srcHandle.LineInfo(29, 16)},
		{token.RBRACKET, "]", srcHandle.LineInfo(29, 17)},
		{token.SEMICOLON, ";", srcHandle.LineInfo(29, 18)},

		{token.EOF, "", srcHandle.LineInfo(30, 0)},
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
	token.ResetForTesting()
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
		FileIndex: 0,
		Line:      3,
		Char:      1,
	}
	if !LineInfoEquals(expectedLineInfo, tok.LineInfo) {
		t.Errorf("expected %s got %s", expectedLineInfo, tok.LineInfo)
	}
}
