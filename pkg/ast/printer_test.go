package ast

import (
	"testing"

	tok "github.com/cedricmar/bazic/pkg/token"
	"github.com/stretchr/testify/assert"
)

func TestPrinter(t *testing.T) {
	expr := NewBinary(
		NewUnary(
			tok.NewToken(tok.MINUS, "-", nil, 1),
			NewLiteral(123),
		),
		tok.NewToken(tok.STAR, "*", nil, 1),
		NewGrouping(
			NewLiteral(45.67),
		),
	)

	assert.Equal(t, "(* (- 123) (group 45.67))", Printer{}.Print(expr))
}
