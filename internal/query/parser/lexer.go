package parser

import "github.com/evanxg852000/foxdb/internal/query/parser/token"

type Lexer struct {
	input string
	curr  int
	peek  int
	ch    byte
}

func NewLexer(input string) *Lexer {
	l := &Lexer{input: input}
	l.consumeChar()
	return l
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token
	l.skipWhitespace()

	switch l.ch {
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case '-':
		tok = newToken(token.MINUS, l.ch)
	case '*':
		tok = newToken(token.ASTERISK, l.ch)
	case '/':
		tok = newToken(token.SLASH, l.ch)
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.consumeChar()
			tok = newToken(token.NOT_EQ, string(ch)+string(l.ch))
		} else {
			tok = newToken(token.BANG, l.ch)
		}
	case '=':
		tok = newToken(token.EQ, l.ch)
	case '<':
		if l.peekChar() == '=' {
			ch := l.ch
			l.consumeChar()
			tok = newToken(token.LT_EQ, string(ch)+string(l.ch))
		} else if l.peekChar() == '>' {
			ch := l.ch
			l.consumeChar()
			tok = newToken(token.NOT_EQ, string(ch)+string(l.ch))
		} else {
			tok = newToken(token.LT, l.ch)
		}
	case '>':
		if l.peekChar() == '=' {
			ch := l.ch
			l.consumeChar()
			tok = newToken(token.GT_EQ, string(ch)+string(l.ch))
		} else {
			tok = newToken(token.GT, l.ch)
		}
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case '.':
		tok = newToken(token.DOT, l.ch)
	case ':':
		tok = newToken(token.COLON, l.ch)
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case '[':
		tok = newToken(token.LBRACKET, l.ch)
	case ']':
		tok = newToken(token.RBRACKET, l.ch)
	case '@':
		tok = newToken(token.AT, l.ch)
	case '#':
		tok = newToken(token.POUND, l.ch)
	case '?':
		tok = newToken(token.QUESTION, l.ch)
	case '&':
		tok = newToken(token.AMPERSAND, l.ch)
	case '"':
		tok.Literal = l.readString()
		tok.Type = token.STRING
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdentifier(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			tok.Literal = l.readNumber()
			tok.Type = token.LookupNumberType(tok.Literal)
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.ch)
		}
	}
	l.consumeChar()
	return tok
}

func (l *Lexer) consumeChar() {
	if l.peek >= len(l.input) {
		l.ch = 0 // eof character
	} else {
		l.ch = l.input[l.peek]
	}
	l.curr = l.peek
	l.peek++
}

func (l *Lexer) peekChar() byte {
	if l.peek >= len(l.input) {
		return 0
	}
	return l.input[l.peek]
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.consumeChar()
	}
}

func (l *Lexer) readString() string {
	l.consumeChar() // skip opening quote
	start := l.curr
	for l.ch != '"' && l.ch != 0 {
		l.consumeChar()
	}
	str := l.input[start:l.curr]
	return str
}

func (l *Lexer) readIdentifier() string {
	start := l.curr
	for isLetter(l.ch) || isDigit(l.ch) {
		l.consumeChar()
	}
	return l.input[start:l.curr]
}

func (l *Lexer) readNumber() string {
	start := l.curr
	for isDigit(l.ch) {
		l.consumeChar()
	}
	if l.ch == '.' {
		l.consumeChar()
		for isDigit(l.ch) {
			l.consumeChar()
		}
	}
	return l.input[start:l.curr]
}

func isLetter(ch byte) bool {
	return ('a' <= ch && ch <= 'z') || ('A' <= ch && ch <= 'Z') || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func newToken[L byte | string](tokenType token.TokenType, literal L) token.Token {
	return token.Token{Type: tokenType, Literal: string(literal)}
}
