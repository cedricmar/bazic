package ast

import (
	"bytes"
	"fmt"
)

// Printer is a "pretty printer" for an AST
type Printer struct{}

func (p Printer) Print(e Expr) string {
	return fmt.Sprintf("%s", e.Accept(p))
}

func (p Printer) VisitBinaryExpr(expr Binary) interface{} {
	return p.parenthesize(expr.operator.Lexeme, expr.left, expr.right)
}

func (p Printer) VisitGroupingExpr(expr Grouping) interface{} {
	return p.parenthesize("group", expr.expression)
}

func (p Printer) VisitLiteralExpr(expr Literal) interface{} {
	if expr.value == nil {
		return "nil"
	}
	return fmt.Sprintf("%v", expr.value)
}

func (p Printer) VisitUnaryExpr(expr Unary) interface{} {
	return p.parenthesize(expr.operator.Lexeme, expr.right)
}

func (p Printer) parenthesize(name string, exprs ...Expr) interface{} {
	var buf bytes.Buffer

	buf.WriteString("(" + name)
	for _, e := range exprs {
		buf.WriteString(" ")
		buf.WriteString(fmt.Sprintf("%s", e.Accept(p)))
	}
	buf.WriteString(")")

	return buf.String()
}
