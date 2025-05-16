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
	reader       *bufio.Reader
	position     int  // Current position in input (points to current char)
	readPosition int  // current reading position in input (after current char)
	ch           rune // current char under examination
}

func NewFromString(input string) *Lexer {
	strReader := strings.NewReader(input)
	r := bufio.NewReader(strReader)
	l := &Lexer{reader: r}
	l.readChar()
	return l
}

func NewFromReader(r io.Reader) *Lexer {
	bufReader := bufio.NewReader(r)
	l := &Lexer{reader: bufReader}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	readRune, _, err := l.reader.ReadRune()
	if errors.Is(err, io.EOF) {
		l.ch = 0
	} else {
		l.ch = readRune
	}
}

func (l *Lexer) NextToken() token.Token {
	l.skipWhitespace()
	readNextChar := true
	var tok token.Token
	switch l.ch {
	case '=':
		peek := l.peekChar()
		if peek == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.EQ, Literal: literal}
		} else {
			tok = newToken(token.ASSIGN, l.ch)
		}
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case '-':
		tok = newToken(token.MINUS, l.ch)
	case '!':
		tok = newToken(token.BANG, l.ch)
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.NOT_EQ, Literal: literal}
		} else {
			tok = newToken(token.BANG, l.ch)
		}
	case '*':
		tok = newToken(token.ASTERISK, l.ch)
	case '/':
		tok = newToken(token.SLASH, l.ch)
	case '<':
		tok = newToken(token.LT, l.ch)
	case '>':
		tok = newToken(token.GT, l.ch)
	case '"':
		tok.Type = token.STRING
		tok.Literal = l.readString()
	case '[':
		tok = newToken(token.LBRACKET, l.ch)
	case ']':
		tok = newToken(token.RBRACKET, l.ch)
	case ':':
		tok = newToken(token.COLON, l.ch)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		readNextChar = false
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			tok.Type = token.INT
			tok.Literal = l.readNumber()
		} else {
			tok = newToken(token.ILLEGAL, l.ch)
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

func (l *Lexer) readNumber() string {
	buffer := make([]rune, 0)
	for isDigit(l.ch) {
		buffer = append(buffer, l.ch)
		l.readChar()
	}
	return string(buffer)
}

func newToken(t token.TokenType, b rune) token.Token {
	return token.Token{
		Type:    t,
		Literal: string(b),
	}
}

func isLetter(b rune) bool {
	return unicode.IsLetter(b)
}

func isDigit(b rune) bool {
	return unicode.IsDigit(b)
}

func (l *Lexer) skipWhitespace() {
	for unicode.IsSpace(l.ch) {
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
