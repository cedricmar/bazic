package parser

import tok "github.com/cedricmar/bazic/pkg/token"

type Parser struct {
	tokens  []tok.Token
	current int
}

func NewParser(tokens []tok.Token) Parser {
	return Parser{tokens, 0}
}
