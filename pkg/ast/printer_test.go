package ast

import (
	"testing"

	"github.com/cedricmar/bazic/pkg/scanner"
	"github.com/stretchr/testify/assert"
)

func TestPrinter(t *testing.T) {
	expr := NewBinary(
		NewUnary(
			scanner.NewToken(scanner.MINUS, "-", nil, 1),
			NewLiteral(123)),
		scanner.NewToken(scanner.STAR, "*", nil, 1),
		NewGrouping(
			NewLiteral(45.67)))

	assert.Equal(t, "(* (- 123) (group 45.67))", Printer{}.Print(expr))
}
