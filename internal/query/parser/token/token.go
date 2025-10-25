package token

import "strings"

type TokenType int

const (
	ILLEGAL TokenType = iota
	EOF

	// Operators
	PLUS     // +
	MINUS    // -
	ASTERISK // *
	SLASH    // /
	BANG     // !

	EQ     // =
	LT     // <
	GT     // >
	LT_EQ  // <=
	GT_EQ  // >=
	NOT_EQ // != or <>

	// Delimiters
	COMMA     // ,
	SEMICOLON // ;
	DOT       // .
	COLON     // :
	LPAREN    // (
	RPAREN    // )
	LBRACE    // {
	RBRACE    // }
	LBRACKET  // [
	RBRACKET  // ]

	// Symbols
	AT        // @
	POUND     // #
	QUESTION  // ?
	AMPERSAND // &

	// Identifiers and literals
	INT    // 12345
	FLOAT  // 123.45
	STRING // "abc"
	IDENT  // main, foo, bar, x, y, z

	// Keywords
	TRUE       // true
	FALSE      // false
	AND        // and
	OR         // or
	PRIMARY    // primary
	KEY        // key
	IF         // if
	NOT        // not
	UNIQUE     // unique
	NULL       // null
	CREATE     // create
	DROP       // drop
	SCHEMA     // schema
	TABLE      // table
	INDEX      // index
	INSERT     // insert
	SELECT     // select
	UPDATE     // update
	DELETE     // delete
	EXISTS     // exists
	INT_TYPE   // int
	FLOAT_TYPE // float
	BOOL_TYPE  // bool
	TEXT_TYPE  // text
)

func (tt TokenType) String() string {
	switch tt {
	case ILLEGAL:
		return "ILLEGAL"
	case EOF:
		return "EOF"
	case PLUS:
		return "+"
	case MINUS:
		return "-"
	case ASTERISK:
		return "*"
	case SLASH:
		return "/"
	case BANG:
		return "!"
	case EQ:
		return "="
	case LT:
		return "<"
	case GT:
		return ">"
	case LT_EQ:
		return "<="
	case GT_EQ:
		return ">="
	case NOT_EQ:
		return "!="
	case COMMA:
		return ","
	case SEMICOLON:
		return ";"
	case DOT:
		return "."
	case COLON:
		return ":"
	case LPAREN:
		return "("
	case RPAREN:
		return ")"
	case LBRACE:
		return "{"
	case RBRACE:
		return "}"
	case LBRACKET:
		return "["
	case RBRACKET:
		return "]"
	case AT:
		return "@"
	case POUND:
		return "#"
	case QUESTION:
		return "?"
	case AMPERSAND:
		return "&"
	case INT:
		return "INT"
	case FLOAT:
		return "FLOAT"
	case STRING:
		return "STRING"
	case IDENT:
		return "IDENT"
	case TRUE:
		return "TRUE"
	case FALSE:
		return "FALSE"
	case NULL:
		return "NULL"
	case SELECT:
		return "SELECT"
	}
	return "UNKNOWN"
}

type Token struct {
	Type    TokenType
	Literal string
}

var keywords = map[string]TokenType{
	"true":   TRUE,
	"false":  FALSE,
	"null":   NULL,
	"create": CREATE,
	"drop":   DROP,
	"table":  TABLE,
	"index":  INDEX,
	"insert": INSERT,
	"select": SELECT,
	"update": UPDATE,
	"delete": DELETE,
}

func LookupIdentifier(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}

func LookupNumberType(literal string) TokenType {
	if strings.Contains(literal, ".") {
		return FLOAT
	}
	return INT
}
