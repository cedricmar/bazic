package token

import "fmt"

type Token struct {
	TokenType TokenType
	Lexeme    string
	Literal   interface{}
	Line      int
}

func NewToken(tokenType TokenType, lexeme string, literal interface{}, line int) Token {
	return Token{tokenType, lexeme, literal, line}
}

func (tok Token) toString() string {
	return fmt.Sprintf("%d %s %v\n", tok.TokenType, tok.Lexeme, tok.Literal)
}
