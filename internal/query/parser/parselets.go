package parser

import (
	"strconv"

	"github.com/evanxg852000/foxdb/internal/query/parser/ast"
	"github.com/evanxg852000/foxdb/internal/query/parser/token"
)

// prefix parselets

func parsePrefixExpression(p *Parser) ast.Expression {
	operator := p.currentToken.Literal
	p.nextToken()

	expression := &ast.PrefixExpr{
		Operator: operator,
		Right:    p.parseExpression(PREFIX),
	}
	return expression
}

func parseIdentifier(p *Parser) ast.Expression {
	return &ast.IdentifierExpr{Value: p.currentToken.Literal}
}

func parseLiteralValue(p *Parser) ast.Expression {
	switch p.currentToken.Type {
	case token.INT:
		value, err := strconv.ParseInt(p.currentToken.Literal, 10, 64)
		if err != nil {
			msg := "could not parse integer literal: " + err.Error()
			p.errors = append(p.errors, msg)
			return nil
		}
		return &ast.IntegerLiteralExpr{Value: value}
	case token.FLOAT:
		value, err := strconv.ParseFloat(p.currentToken.Literal, 64)
		if err != nil {
			msg := "could not parse float literal: " + err.Error()
			p.errors = append(p.errors, msg)
			return nil
		}
		return &ast.FloatLiteralExpr{Value: value}
	case token.STRING:
		return &ast.StringLiteralExpr{Value: p.currentToken.Literal}
	case token.NULL:
		return &ast.NullLiteralExpr{}
	case token.TRUE:
		return &ast.BooleanLiteralExpr{Value: true}
	case token.FALSE:
		return &ast.BooleanLiteralExpr{Value: false}
	}
	return nil
}

func parseGroupedExpression(p *Parser) ast.Expression {
	p.nextToken()

	exp := p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	return exp
}

// infix parselets

func parseInfixExpression(p *Parser, left ast.Expression) ast.Expression {
	operator := p.currentToken.Literal
	precedence := p.currentPrecedence()

	p.nextToken()

	expression := &ast.InfixExpr{
		Left:     left,
		Operator: operator,
		Right:    p.parseExpression(precedence),
	}
	return expression
}

func parseCallExpression(p *Parser, function ast.Expression) ast.Expression {
	exp := &ast.CallExpr{Function: function}
	exp.Args = p.parseExpressionList(token.RPAREN)
	return exp
}
