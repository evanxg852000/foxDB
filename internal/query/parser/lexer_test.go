package parser

import (
	"testing"

	"github.com/evanxg852000/foxdb/internal/query/parser/token"
	"github.com/stretchr/testify/assert"
)

func TestLexerNextToken(t *testing.T) {
	input := `=+-*/!< > (){}[],.;:@#?&`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.EQ, "="},
		{token.PLUS, "+"},
		{token.MINUS, "-"},
		{token.ASTERISK, "*"},
		{token.SLASH, "/"},
		{token.BANG, "!"},
		{token.LT, "<"},
		{token.GT, ">"},
		{token.LPAREN, "("},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RBRACE, "}"},
		{token.LBRACKET, "["},
		{token.RBRACKET, "]"},
		{token.COMMA, ","},
		{token.DOT, "."},
		{token.SEMICOLON, ";"},
		{token.COLON, ":"},
		{token.AT, "@"},
		{token.POUND, "#"},
		{token.QUESTION, "?"},
		{token.AMPERSAND, "&"},
		{token.EOF, ""},
	}

	l := NewLexer(input)

	for i, tt := range tests {
		tok := l.NextToken()
		assert.Equal(t, tt.expectedType, tok.Type, "test[%d] - unexpected token type", i)
		assert.Equal(t, tt.expectedLiteral, tok.Literal, "test[%d] - unexpected token literal", i)
	}
}

func TestLexerComparisonOperators(t *testing.T) {
	input := `<= >= != <>`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.LT_EQ, "<="},
		{token.GT_EQ, ">="},
		{token.NOT_EQ, "!="},
		{token.NOT_EQ, "<>"},
		{token.EOF, ""},
	}

	l := NewLexer(input)

	for i, tt := range tests {
		tok := l.NextToken()
		assert.Equal(t, tt.expectedType, tok.Type, "test[%d] - unexpected token type", i)
		assert.Equal(t, tt.expectedLiteral, tok.Literal, "test[%d] - unexpected token literal", i)
	}
}

func TestLexerIdentifiers(t *testing.T) {
	input := `foo bar_baz myVariable _underscore ABC123`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.IDENT, "foo"},
		{token.IDENT, "bar_baz"},
		{token.IDENT, "myVariable"},
		{token.IDENT, "_underscore"},
		{token.IDENT, "ABC123"},
		{token.EOF, ""},
	}

	l := NewLexer(input)

	for i, tt := range tests {
		tok := l.NextToken()
		assert.Equal(t, tt.expectedType, tok.Type, "test[%d] - unexpected token type", i)
		assert.Equal(t, tt.expectedLiteral, tok.Literal, "test[%d] - unexpected token literal", i)
	}
}

func TestLexerKeywords(t *testing.T) {
	input := `true false null select`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.TRUE, "true"},
		{token.FALSE, "false"},
		{token.NULL, "null"},
		{token.SELECT, "select"},
		{token.EOF, ""},
	}

	l := NewLexer(input)

	for i, tt := range tests {
		tok := l.NextToken()
		assert.Equal(t, tt.expectedType, tok.Type, "test[%d] - unexpected token type", i)
		assert.Equal(t, tt.expectedLiteral, tok.Literal, "test[%d] - unexpected token literal", i)
	}
}

func TestLexerNumbers(t *testing.T) {
	input := `123 456.789 0 0.0 999.99`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.INT, "123"},
		{token.FLOAT, "456.789"},
		{token.INT, "0"},
		{token.FLOAT, "0.0"},
		{token.FLOAT, "999.99"},
		{token.EOF, ""},
	}

	l := NewLexer(input)

	for i, tt := range tests {
		tok := l.NextToken()
		assert.Equal(t, tt.expectedType, tok.Type, "test[%d] - unexpected token type", i)
		assert.Equal(t, tt.expectedLiteral, tok.Literal, "test[%d] - unexpected token literal", i)
	}
}

func TestLexerStrings(t *testing.T) {
	input := `"hello" "world with spaces" "special !@#$ chars" ""`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.STRING, "hello"},
		{token.STRING, "world with spaces"},
		{token.STRING, "special !@#$ chars"},
		{token.STRING, ""},
		{token.EOF, ""},
	}

	l := NewLexer(input)

	for i, tt := range tests {
		tok := l.NextToken()
		assert.Equal(t, tt.expectedType, tok.Type, "test[%d] - unexpected token type", i)
		assert.Equal(t, tt.expectedLiteral, tok.Literal, "test[%d] - unexpected token literal", i)
	}
}

func TestLexerWhitespace(t *testing.T) {
	input := `  foo   bar	
	baz  `

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.IDENT, "foo"},
		{token.IDENT, "bar"},
		{token.IDENT, "baz"},
		{token.EOF, ""},
	}

	l := NewLexer(input)

	for i, tt := range tests {
		tok := l.NextToken()
		assert.Equal(t, tt.expectedType, tok.Type, "test[%d] - unexpected token type", i)
		assert.Equal(t, tt.expectedLiteral, tok.Literal, "test[%d] - unexpected token literal", i)
	}
}

func TestLexerIllegalCharacters(t *testing.T) {
	input := `$ % ^`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.ILLEGAL, "$"},
		{token.ILLEGAL, "%"},
		{token.ILLEGAL, "^"},
		{token.EOF, ""},
	}

	l := NewLexer(input)

	for i, tt := range tests {
		tok := l.NextToken()
		assert.Equal(t, tt.expectedType, tok.Type, "test[%d] - unexpected token type", i)
		assert.Equal(t, tt.expectedLiteral, tok.Literal, "test[%d] - unexpected token literal", i)
	}
}

func TestLexerComplexSQL(t *testing.T) {
	input := `select * FROM users WHERE age >= 18 AND name != "admin"`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.SELECT, "select"},
		{token.ASTERISK, "*"},
		{token.IDENT, "FROM"},
		{token.IDENT, "users"},
		{token.IDENT, "WHERE"},
		{token.IDENT, "age"},
		{token.GT_EQ, ">="},
		{token.INT, "18"},
		{token.IDENT, "AND"},
		{token.IDENT, "name"},
		{token.NOT_EQ, "!="},
		{token.STRING, "admin"},
		{token.EOF, ""},
	}

	l := NewLexer(input)

	for i, tt := range tests {
		tok := l.NextToken()
		assert.Equal(t, tt.expectedType, tok.Type, "test[%d] - unexpected token type", i)
		assert.Equal(t, tt.expectedLiteral, tok.Literal, "test[%d] - unexpected token literal", i)
	}
}

func TestLexerEmptyString(t *testing.T) {
	input := ""
	l := NewLexer(input)

	tok := l.NextToken()
	assert.Equal(t, token.EOF, tok.Type, "expected EOF for empty input")
	assert.Equal(t, "", tok.Literal, "expected empty literal for EOF")
}

func TestLexerSingleCharacter(t *testing.T) {
	tests := []struct {
		input           string
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{"+", token.PLUS, "+"},
		{"-", token.MINUS, "-"},
		{"*", token.ASTERISK, "*"},
		{"/", token.SLASH, "/"},
		{"!", token.BANG, "!"},
		{"=", token.EQ, "="},
		{"<", token.LT, "<"},
		{">", token.GT, ">"},
		{"(", token.LPAREN, "("},
		{")", token.RPAREN, ")"},
		{"{", token.LBRACE, "{"},
		{"}", token.RBRACE, "}"},
		{"[", token.LBRACKET, "["},
		{"]", token.RBRACKET, "]"},
		{",", token.COMMA, ","},
		{".", token.DOT, "."},
		{";", token.SEMICOLON, ";"},
		{":", token.COLON, ":"},
		{"@", token.AT, "@"},
		{"#", token.POUND, "#"},
		{"?", token.QUESTION, "?"},
		{"&", token.AMPERSAND, "&"},
	}

	for _, tt := range tests {
		l := NewLexer(tt.input)
		tok := l.NextToken()

		assert.Equal(t, tt.expectedType, tok.Type, "input %q - unexpected token type", tt.input)
		assert.Equal(t, tt.expectedLiteral, tok.Literal, "input %q - unexpected token literal", tt.input)

		// Should return EOF for next token
		tok = l.NextToken()
		assert.Equal(t, token.EOF, tok.Type, "input %q - expected EOF after single token", tt.input)
	}
}

func TestLexerUnterminatedString(t *testing.T) {
	input := `"unterminated string`
	l := NewLexer(input)

	tok := l.NextToken()
	assert.Equal(t, token.STRING, tok.Type, "expected STRING token for unterminated string")
	assert.Equal(t, "unterminated string", tok.Literal, "unexpected literal for unterminated string")
}

func TestLexerMixedTokens(t *testing.T) {
	input := ` foo123 = "bar" + 456.78; ( true != false )`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.IDENT, "foo123"},
		{token.EQ, "="},
		{token.STRING, "bar"},
		{token.PLUS, "+"},
		{token.FLOAT, "456.78"},
		{token.SEMICOLON, ";"},
		{token.LPAREN, "("},
		{token.TRUE, "true"},
		{token.NOT_EQ, "!="},
		{token.FALSE, "false"},
		{token.RPAREN, ")"},
		{token.EOF, ""},
	}

	l := NewLexer(input)

	for i, tt := range tests {
		tok := l.NextToken()
		assert.Equal(t, tt.expectedType.String(), tok.Type.String(), "test[%d] - unexpected token type", i)
		assert.Equal(t, tt.expectedLiteral, tok.Literal, "test[%d] - unexpected token literal", i)
	}
}

func TestLookupIdentifier(t *testing.T) {
	tests := []struct {
		input    string
		expected token.TokenType
	}{
		{"true", token.TRUE},
		{"false", token.FALSE},
		{"null", token.NULL},
		{"select", token.SELECT},
		{"foo", token.IDENT},
		{"bar", token.IDENT},
		{"variable_name", token.IDENT},
		{"TRUE", token.IDENT},   // case sensitive
		{"SELECT", token.IDENT}, // case sensitive
	}

	for _, tt := range tests {
		result := token.LookupIdentifier(tt.input)
		assert.Equal(t, tt.expected, result, "LookupIdentifier(%q) returned unexpected result", tt.input)
	}
}

func TestLookupNumberType(t *testing.T) {
	tests := []struct {
		input    string
		expected token.TokenType
	}{
		{"123", token.INT},
		{"0", token.INT},
		{"999", token.INT},
		{"123.456", token.FLOAT},
		{"0.0", token.FLOAT},
		{"999.99", token.FLOAT},
		{"1.0", token.FLOAT},
	}

	for _, tt := range tests {
		result := token.LookupNumberType(tt.input)
		assert.Equal(t, tt.expected, result, "LookupNumberType(%q) returned unexpected result", tt.input)
	}
}
