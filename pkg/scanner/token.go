package scanner

import "fmt"

type Token struct {
	tokenType TokenType
	Lexeme    string
	literal   interface{}
	line      int
}

func NewToken(tokenType TokenType, lexeme string, literal interface{}, line int) Token {
	return Token{tokenType, lexeme, literal, line}
}

func (t Token) toString() string {
	return fmt.Sprintf("%d %s %v\n", t.tokenType, t.Lexeme, t.literal)
}
