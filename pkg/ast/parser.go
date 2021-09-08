package ast

import (
	"github.com/cedricmar/bazic/pkg/scanner"

	tok "github.com/cedricmar/bazic/pkg/token"
)

// Parser uses Recursive Descent Parsing
type Parser struct {
	tokens  []tok.Token
	current int
}

type ParseError struct {
	token tok.Token
	msg   string
}

func NewParser(tokens []tok.Token) Parser {
	return Parser{tokens, 0}
}

func (p *Parser) Parse() (Expr, error) {
	return p.Expression()
}

func (e ParseError) Error() string {
	scanner.Error(e.token, e.msg)
	return ""
}

// expression     → equality
func (p *Parser) Expression() (Expr, error) {
	return p.Equality()
}

// equality       → comparison ( ( "!=" | "==" ) comparison )*
func (p *Parser) Equality() (Expr, error) {
	expr, err := p.Comparison()
	if err != nil {
		return expr, err
	}

	for p.match(tok.BANG_EQUAL, tok.EQUAL_EQUAL) {
		operator := p.previous()
		right, err := p.Comparison()
		if err != nil {
			return right, err
		}
		expr = NewBinary(expr, operator, right)

	}

	return expr, nil
}

// comparison     → term ( ( ">" | ">=" | "<" | "<=" ) term )*
func (p *Parser) Comparison() (Expr, error) {
	expr, err := p.term()
	if err != nil {
		return expr, err
	}

	for p.match(tok.GREATER, tok.GREATER_EQUAL, tok.LESS, tok.LESS_EQUAL) {
		operator := p.previous()
		right, err := p.term()
		if err != nil {
			return right, err
		}
		expr = NewBinary(expr, operator, right)
	}

	return expr, nil
}

// unary          → ( "!" | "-" ) unary | primary
func (p *Parser) Unary() (Expr, error) {

	if p.match(tok.BANG, tok.MINUS) {
		operator := p.previous()
		right, err := p.Unary()
		if err != nil {
			return right, err
		}
		return NewUnary(operator, right), nil
	}

	return p.Primary()
}

// primary        → NUMBER | STRING | "true" | "false" | "nil" | "(" expression ")"
func (p *Parser) Primary() (Expr, error) {
	if p.match(tok.FALSE) {
		return NewLiteral(false), nil
	}
	if p.match(tok.TRUE) {
		return NewLiteral(true), nil
	}
	if p.match(tok.NIL) {
		return NewLiteral(nil), nil
	}

	if p.match(tok.NUMBER, tok.STRING) {
		return NewLiteral(p.previous().Literal), nil
	}

	if p.match(tok.LEFT_PAREN) {
		expr, err := p.Expression()
		if err != nil {
			return expr, err
		}
		p.consume(tok.RIGHT_PAREN, "Expect ')' after expression.")
		return NewGrouping(expr), nil
	}

	return Literal{}, ParseError{
		token: p.peek(),
		msg:   "Expect expression.",
	}
}

func (p *Parser) consume(t tok.TokenType, m string) (tok.Token, error) {
	if p.check(t) {
		return p.advance(), nil
	}

	return tok.Token{}, ParseError{
		token: p.peek(),
		msg:   m,
	}
}

func (p *Parser) term() (Expr, error) {
	expr, err := p.factor()
	if err != nil {
		return expr, err
	}

	for p.match(tok.MINUS, tok.PLUS) {
		operator := p.previous()
		right, err := p.factor()
		if err != nil {
			return right, err
		}
		expr = NewBinary(expr, operator, right)
	}

	return expr, nil
}

func (p *Parser) factor() (Expr, error) {
	expr, err := p.Unary()
	if err != nil {
		return expr, err
	}

	for p.match(tok.SLASH, tok.STAR) {
		operator := p.previous()
		right, err := p.Unary()
		if err != nil {
			return right, err
		}
		expr = NewBinary(expr, operator, right)
	}

	return expr, nil
}

func (p *Parser) match(types ...tok.TokenType) bool {
	for _, t := range types {
		if p.check(t) {
			p.advance()
			return true
		}
	}

	return false
}

func (p *Parser) check(t tok.TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().TokenType == t
}

func (p *Parser) advance() tok.Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) isAtEnd() bool {
	return p.peek().TokenType == tok.EOF
}

func (p *Parser) peek() tok.Token {
	return p.tokens[p.current]
}

func (p *Parser) previous() tok.Token {
	return p.tokens[p.current-1]
}

func (p *Parser) synchronize() {
	p.advance()

	for !p.isAtEnd() {

		if p.previous().TokenType == tok.SEMICOLON {
			return
		}

		switch p.peek().TokenType {
		case tok.CLASS:
		case tok.FUN:
		case tok.VAR:
		case tok.FOR:
		case tok.IF:
		case tok.WHILE:
		case tok.PRINT:
		case tok.RETURN:
			return
		}

		p.advance()
	}
}
