package lexer

import (
	"bufio"
	"errors"
	"io"
	"monkey/token"
	"strings"
	"unicode"
)

type Lexer struct {
	reader   *bufio.Reader
	fileName string
	lineNo   int
	charNo   int
	ch       rune // current char under examination
}

func NewFromString(name, input string) *Lexer {
	strReader := strings.NewReader(input)
	r := bufio.NewReader(strReader)
	l := &Lexer{reader: r, fileName: name, lineNo: 1, charNo: 0}
	l.readChar()
	return l
}

func NewFromReader(fileName string, r io.Reader) *Lexer {
	bufReader := bufio.NewReader(r)
	l := &Lexer{reader: bufReader, fileName: fileName, lineNo: 1, charNo: 0}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	readRune, _, err := l.reader.ReadRune()
	if errors.Is(err, io.EOF) {
		l.ch = 0
	} else {
		l.charNo++
		l.ch = readRune
	}
}

func (l *Lexer) NextToken() token.Token {
	l.skipWhitespace()
	readNextChar := true
	var tok token.Token
	lineInfo := token.LineInfo{FileName: l.fileName, Line: l.lineNo, Char: l.charNo}
	switch l.ch {
	case '=':
		peek := l.peekChar()
		if peek == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.EQ, Literal: literal, LineInfo: lineInfo}
		} else {
			tok = newToken(token.ASSIGN, l.ch, lineInfo)
		}
	case ';':
		tok = newToken(token.SEMICOLON, l.ch, lineInfo)
	case '(':
		tok = newToken(token.LPAREN, l.ch, lineInfo)
	case ')':
		tok = newToken(token.RPAREN, l.ch, lineInfo)
	case ',':
		tok = newToken(token.COMMA, l.ch, lineInfo)
	case '+':
		tok = newToken(token.PLUS, l.ch, lineInfo)
	case '{':
		tok = newToken(token.LBRACE, l.ch, lineInfo)
	case '}':
		tok = newToken(token.RBRACE, l.ch, lineInfo)
	case '-':
		tok = newToken(token.MINUS, l.ch, lineInfo)
	case '!':
		tok = newToken(token.BANG, l.ch, lineInfo)
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.NOT_EQ, Literal: literal, LineInfo: lineInfo}
		} else {
			tok = newToken(token.BANG, l.ch, lineInfo)
		}
	case '*':
		tok = newToken(token.ASTERISK, l.ch, lineInfo)
	case '/':
		tok = newToken(token.SLASH, l.ch, lineInfo)
	case '<':
		tok = newToken(token.LT, l.ch, lineInfo)
	case '>':
		tok = newToken(token.GT, l.ch, lineInfo)
	case '"':
		tok.Type = token.STRING
		tok.Literal = l.readString()
		tok.LineInfo = lineInfo
	case '[':
		tok = newToken(token.LBRACKET, l.ch, lineInfo)
	case ']':
		tok = newToken(token.RBRACKET, l.ch, lineInfo)
	case ':':
		tok = newToken(token.COLON, l.ch, lineInfo)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
		tok.LineInfo = lineInfo
	default:
		readNextChar = false
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			tok.LineInfo = lineInfo
			return tok
		} else if isDigit(l.ch) {
			tok.Type, tok.Literal = l.readNumber()
			tok.LineInfo = lineInfo
		} else {
			tok = newToken(token.ILLEGAL, l.ch, lineInfo)
		}
	}

	if readNextChar {
		l.readChar()
	}
	return tok
}

func (l *Lexer) peekChar() rune {
	readRune, _, err := l.reader.ReadRune()
	if err != nil {
		return 0
	}
	l.reader.UnreadRune()
	return readRune
}

func (l *Lexer) readIdentifier() string {
	buffer := make([]rune, 0)
	for isLetter(l.ch) {
		buffer = append(buffer, l.ch)
		l.readChar()
	}
	return string(buffer)
}

// Octal numbers start with 0
// Hex numbers 0x
// floating point numbers
func (l *Lexer) readNumber() (token.TokenType, string) {
	numType := token.INT
	buffer := make([]rune, 0)
	classifier := isDigit
	if l.ch == '0' {
		buffer = append(buffer, l.ch)
		l.readChar()
		if l.ch == 'x' || l.ch == 'X' {
			classifier = isHexDigit
		} else if isOctalDigit(l.ch) {
			classifier = isOctalDigit
		}
	}

	for l.ch == '.' || classifier(l.ch) {
		if l.ch == '.' {
			numType = token.FLOAT
		}
		buffer = append(buffer, l.ch)
		l.readChar()
	}
	return token.TokenType(numType), string(buffer)
}

func newToken(t token.TokenType, b rune, lineInfo token.LineInfo) token.Token {
	return token.Token{
		Type:     t,
		Literal:  string(b),
		LineInfo: lineInfo,
	}
}

func isLetter(b rune) bool {
	return unicode.IsLetter(b)
}
func isHexDigit(b rune) bool {
	return isDigit(b) || (b >= 65 && b <= 90) || (b >= 97 && b <= 122)
}

func isOctalDigit(b rune) bool {
	return b >= 48 && b >= 55
}

func isDigit(b rune) bool {
	return unicode.IsDigit(b)
}

func (l *Lexer) skipWhitespace() {
	for unicode.IsSpace(l.ch) {
		if l.ch == '\n' {
			l.charNo = 0
			l.lineNo++
		}
		l.readChar()
	}
}

func (l *Lexer) readString() string {
	buffer := make([]rune, 0)
	for {
		l.readChar()
		if l.ch == '"' || l.ch == 0 {
			break
		}
		buffer = append(buffer, l.ch)
	}
	return string(buffer)
}
