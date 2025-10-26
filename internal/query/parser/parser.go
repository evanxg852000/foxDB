package parser

import (
	"fmt"

	"github.com/evanxg852000/foxdb/internal/query/parser/ast"
	"github.com/evanxg852000/foxdb/internal/query/parser/token"
	"github.com/evanxg852000/foxdb/internal/types"
)

// operator precedences
const (
	_ int = iota
	LOWEST
	AND_OR      // AND, OR
	COMP        // ==, !=, <, >=, >, <=
	SUM         // +, -
	PRODUCT     // *, /
	PREFIX      // -x, !x
	CALL        // fn(x)
	ARRAY_INDEX // arr[i]
)

var precedencesTable = map[token.TokenType]int{
	token.AND:      AND_OR,
	token.OR:       AND_OR,
	token.EQ:       COMP,
	token.NOT_EQ:   COMP,
	token.LT:       COMP,
	token.LT_EQ:    COMP,
	token.GT:       COMP,
	token.GT_EQ:    COMP,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.ASTERISK: PRODUCT,
	token.SLASH:    PRODUCT,
}

type prefixParseFn func(p *Parser) ast.Expression
type infixParseFn func(p *Parser, left ast.Expression) ast.Expression

type Parser struct {
	lexer          *Lexer
	errors         []string
	currentToken   token.Token
	peekToken      token.Token
	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

func NewParser(lexer *Lexer) *Parser {
	parser := &Parser{
		lexer:          lexer,
		errors:         []string{},
		prefixParseFns: make(map[token.TokenType]prefixParseFn),
		infixParseFns:  make(map[token.TokenType]infixParseFn),
	}

	// Initialize the maps
	parser.prefixParseFns[token.IDENT] = parseIdentifier
	parser.prefixParseFns[token.INT] = parseLiteralValue
	parser.prefixParseFns[token.FLOAT] = parseLiteralValue
	parser.prefixParseFns[token.STRING] = parseLiteralValue
	parser.prefixParseFns[token.NULL] = parseLiteralValue
	parser.prefixParseFns[token.TRUE] = parseLiteralValue
	parser.prefixParseFns[token.FALSE] = parseLiteralValue
	parser.prefixParseFns[token.MINUS] = parsePrefixExpression
	parser.prefixParseFns[token.NOT] = parsePrefixExpression
	parser.prefixParseFns[token.LPAREN] = parseGroupedExpression

	parser.infixParseFns[token.PLUS] = parseInfixExpression
	parser.infixParseFns[token.MINUS] = parseInfixExpression
	parser.infixParseFns[token.ASTERISK] = parseInfixExpression
	parser.infixParseFns[token.SLASH] = parseInfixExpression
	parser.infixParseFns[token.EQ] = parseInfixExpression
	parser.infixParseFns[token.NOT_EQ] = parseInfixExpression
	parser.infixParseFns[token.LT] = parseInfixExpression
	parser.infixParseFns[token.LT_EQ] = parseInfixExpression
	parser.infixParseFns[token.GT] = parseInfixExpression
	parser.infixParseFns[token.GT_EQ] = parseInfixExpression
	parser.infixParseFns[token.AND] = parseInfixExpression
	parser.infixParseFns[token.OR] = parseInfixExpression
	parser.infixParseFns[token.LPAREN] = parseCallExpression

	// Read two tokens, so currentToken and peekToken are both set
	parser.nextToken()
	parser.nextToken()

	return parser
}

func (p *Parser) nextToken() {
	p.currentToken = p.peekToken
	p.peekToken = p.lexer.NextToken()
}

func (p *Parser) currentTokenIs(t token.TokenType) bool {
	return p.currentToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekTokenError(t)
		return false
	}
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) currentTokenError(t token.TokenType) {
	msg := fmt.Sprintf("expected current token to be %s, got %s instead", t, p.currentToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) peekTokenError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{
		Statements: []ast.Statement{},
	}
	for !p.currentTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}
	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.currentToken.Type {
	case token.CREATE:
		return p.parseCreateStatement()
	case token.DROP:
		return p.parseDropStatement()
	// case token.INSERT:
	// 	return p.parseInsertStatement()
	// case token.SELECT:
	// 	return p.parseSelectStatement()
	// case token.UPDATE:
	// 	return p.parseUpdateStatement()
	// case token.DELETE:
	// 	return p.parseDeleteStatement()
	default:
		p.errors = append(p.errors, fmt.Sprintf("unknown statement: %s", p.currentToken.Type))
		return nil
	}
}

func (p *Parser) parseCreateStatement() ast.Statement {
	p.nextToken() // consume 'CREATE'
	switch p.currentToken.Type {
	case token.SCHEMA:
		return p.parseCreateSchemaStatement()
	case token.TABLE:
		return p.parseCreateTableStatement()
	case token.INDEX:
		return p.parseCreateIndexStatement()
	default:
		p.errors = append(p.errors, fmt.Sprintf("expected SCHEMA, TABLE, or INDEX after CREATE, got %s instead", p.currentToken.Type))
		return nil
	}
}

func (p *Parser) parseCreateSchemaStatement() ast.Statement {
	p.nextToken() // consume 'SCHEMA'

	IfNotExists := false
	if p.currentToken.Type == token.IF {
		p.nextToken() // consume 'IF'

		if !p.currentTokenIs(token.NOT) {
			p.currentTokenError(token.NOT)
			return nil
		}
		p.nextToken() // consume 'NOT'

		if !p.currentTokenIs(token.EXISTS) {
			p.currentTokenError(token.EXISTS)
			return nil
		}
		p.nextToken() // consume 'EXISTS'

		IfNotExists = true
	}

	if !p.currentTokenIs(token.IDENT) {
		p.currentTokenError(token.IDENT)
		return nil
	}
	schemaName := p.currentToken.Literal
	p.nextToken() // consume schema name

	if !p.currentTokenIs(token.SEMICOLON) {
		p.currentTokenError(token.SEMICOLON)
		return nil
	}

	return &ast.CreateSchemaStatement{
		SchemaName:  schemaName,
		IfNotExists: IfNotExists,
	}
}

func (p *Parser) parseCreateTableStatement() ast.Statement {
	p.nextToken() // consume 'TABLE'

	if p.currentToken.Type != token.IDENT {
		p.errors = append(p.errors, fmt.Sprintf("expected table name after CREATE TABLE, got %s instead", p.currentToken.Type))
		return nil
	}
	tableName := p.currentToken.Literal
	p.nextToken() // consume table name

	if !p.currentTokenIs(token.LPAREN) {
		p.errors = append(p.errors, fmt.Sprintf("expected '(' after table name, got %s instead", p.currentToken.Type))
		return nil
	}

	columns := []ast.ColumnDef{}

	for p.currentToken.Type != token.RPAREN {
		if p.currentToken.Type != token.IDENT {
			p.errors = append(p.errors, fmt.Sprintf("expected column name, got %s instead", p.currentToken.Type))
			return nil
		}
		columnName := p.currentToken.Literal
		p.nextToken() // consume column name

		var dataType types.DataType
		switch p.currentToken.Type {
		case token.INT_TYPE:
			dataType = types.TYPE_INT
		case token.FLOAT_TYPE:
			dataType = types.TYPE_FLOAT
		case token.BOOL_TYPE:
			dataType = types.TYPE_BOOL
		case token.TEXT_TYPE:
			dataType = types.TYPE_TEXT
		default:
			p.errors = append(p.errors, fmt.Sprintf("expected data type for column %s, got %s instead", columnName, p.currentToken.Type))
			return nil
		}
		p.nextToken() // consume data type

		constraint := ast.Constraint{}
		for p.currentToken.Type != token.COMMA && p.currentToken.Type != token.RPAREN {
			if p.currentToken.Type == token.PRIMARY {
				p.nextToken() // consume PRIMARY
				if p.currentToken.Type != token.KEY {
					p.errors = append(p.errors, fmt.Sprintf("expected KEY after PRIMARY, got %s instead", p.currentToken.Type))
					return nil
				}
				constraint.PrimaryKey = true
				p.nextToken() // consume KEY
			}

			if p.currentToken.Type == token.NOT {
				p.nextToken() // consume NOT
				if p.currentToken.Type != token.NULL {
					p.errors = append(p.errors, fmt.Sprintf("expected NULL after NOT, got %s instead", p.currentToken.Type))
					return nil
				}
				constraint.NotNull = true
				p.nextToken() // consume NULL
			}

			if p.currentToken.Type == token.UNIQUE {
				constraint.Unique = true
				p.nextToken() // consume UNIQUE
			}
		}

		columnDef := ast.ColumnDef{
			Name:       columnName,
			DataType:   dataType,
			Constraint: constraint,
		}
		columns = append(columns, columnDef)

		if p.currentToken.Type == token.COMMA {
			p.nextToken() // consume comma and continue to next column
		} else if p.currentToken.Type != token.RPAREN {
			p.errors = append(p.errors, fmt.Sprintf("expected comma or closing parenthesis, got %s instead", p.currentToken.Type))
			return nil
		}
	}

	if !p.currentTokenIs(token.RPAREN) {
		p.errors = append(p.errors, fmt.Sprintf("expected closing parenthesis, got %s instead", p.currentToken.Type))
		return nil
	}

	if !p.currentTokenIs(token.SEMICOLON) {
		p.errors = append(p.errors, fmt.Sprintf("expected SEMICOLON after CREATE TABLE, got %s instead", p.currentToken.Type))
		return nil
	}

	return &ast.CreateTableStatement{
		TableName: tableName,
		Columns:   columns,
	}
}

func (p *Parser) parseCreateIndexStatement() ast.Statement {
	panic("CREATE INDEX not implemented yet")
}

func (p *Parser) parseDropStatement() ast.Statement {
	p.nextToken() // consume 'DROP'
	switch p.currentToken.Type {
	case token.TABLE:
		return p.parseDropTableStatement()
	case token.INDEX:
		return p.parseDropIndexStatement()
	default:
		p.errors = append(p.errors, fmt.Sprintf("expected TABLE after DROP, got %s instead", p.currentToken.Type))
		return nil
	}
}

func (p *Parser) parseDropTableStatement() ast.Statement {
	p.nextToken() // consume 'TABLE'
	if !p.currentTokenIs(token.IDENT) {
		p.errors = append(p.errors, fmt.Sprintf("expected table name after DROP TABLE, got %s instead", p.currentToken.Type))
		return nil
	}

	if !p.currentTokenIs(token.IDENT) {
		p.errors = append(p.errors, fmt.Sprintf("expected table name after DROP TABLE, got %s instead", p.currentToken.Type))
		return nil
	}
	tableName := p.currentToken.Literal

	p.nextToken() // consume table name

	if !p.currentTokenIs(token.SEMICOLON) {
		p.errors = append(p.errors, fmt.Sprintf("expected semicolon after DROP TABLE, got %s instead", p.currentToken.Type))
		return nil
	}

	return &ast.DropTableStatement{
		TableName: tableName,
	}
}

func (p *Parser) parseDropIndexStatement() ast.Statement {
	panic("DROP INDEX not implemented yet")
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.currentToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.currentToken.Type)
		return nil
	}
	leftExp := prefix(p)

	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}

		p.nextToken()

		leftExp = infix(p, leftExp)
	}

	return leftExp
}

func (p *Parser) parseExpressionList(stopToken token.TokenType) []ast.Expression {
	list := []ast.Expression{}

	if p.peekTokenIs(stopToken) {
		p.nextToken()
		return list
	}

	p.nextToken()
	list = append(list, p.parseExpression(LOWEST))

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		list = append(list, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(stopToken) {
		return nil
	}

	return list
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedencesTable[p.peekToken.Type]; ok {
		return p
	}

	return LOWEST
}

func (p *Parser) currentPrecedence() int {
	if p, ok := precedencesTable[p.currentToken.Type]; ok {
		return p
	}

	return LOWEST
}
