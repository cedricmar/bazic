package ast

import "github.com/cedricmar/bazic/pkg/scanner"

// Expr is a type for the AST
type Expr string

// Binary is a node of the AST
type Binary struct {
	left     Expr
	operator scanner.Token
	right    Expr
}

// NewBinary returns a new node of type Binary
func NewBinary(left Expr, operator scanner.Token, right Expr) Binary {
	return Binary{
		left:     left,
		operator: operator,
		right:    right,
	}
}

// Grouping is a node of the AST
type Grouping struct {
	expression Expr
}

// NewGrouping returns a new node of type Grouping
func NewGrouping(expression Expr) Grouping {
	return Grouping{
		expression: expression,
	}
}

// Literal is a node of the AST
type Literal struct {
	value interface{}
}

// NewLiteral returns a new node of type Literal
func NewLiteral(value interface{}) Literal {
	return Literal{
		value: value,
	}
}

// Unary is a node of the AST
type Unary struct {
	operator scanner.Token
	right    Expr
}

// NewUnary returns a new node of type Unary
func NewUnary(operator scanner.Token, right Expr) Unary {
	return Unary{
		operator: operator,
		right:    right,
	}
}
