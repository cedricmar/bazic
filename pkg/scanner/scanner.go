package scanner

import (
	"fmt"
	"strconv"

	tok "github.com/cedricmar/bazic/pkg/token"
)

type Scanner struct {
	source   string
	tokens   []tok.Token
	start    int
	current  int
	line     int
	HadError bool
}

var keywords = map[string]tok.TokenType{
	"and":    tok.AND,
	"class":  tok.CLASS,
	"else":   tok.ELSE,
	"false":  tok.FALSE,
	"for":    tok.FOR,
	"fun":    tok.FUN,
	"if":     tok.IF,
	"nil":    tok.NIL,
	"or":     tok.OR,
	"print":  tok.PRINT,
	"return": tok.RETURN,
	"super":  tok.SUPER,
	"this":   tok.THIS,
	"true":   tok.TRUE,
	"var":    tok.VAR,
	"while":  tok.WHILE,
}

func NewScanner(source string) Scanner {
	return Scanner{source, []tok.Token{}, 0, 0, 1, false}
}

func (s *Scanner) ScanTokens() []tok.Token {
	for !s.isAtEnd() {
		s.start = s.current
		s.scanToken()
	}

	s.tokens = append(s.tokens, tok.Token{tok.EOF, "", "", s.line})
	return s.tokens
}

func (s *Scanner) scanToken() {
	c := string(s.advance())
	switch c {
	case "(":
		s.addToken(tok.LEFT_PAREN)
		break
	case ")":
		s.addToken(tok.RIGHT_PAREN)
		break
	case "{":
		s.addToken(tok.LEFT_BRACE)
		break
	case "}":
		s.addToken(tok.RIGHT_BRACE)
		break
	case ",":
		s.addToken(tok.COMMA)
		break
	case ".":
		s.addToken(tok.DOT)
		break
	case "-":
		s.addToken(tok.MINUS)
		break
	case "+":
		s.addToken(tok.PLUS)
		break
	case ";":
		s.addToken(tok.SEMICOLON)
		break
	case "*":
		s.addToken(tok.STAR)
		break
	case "!":
		b := tok.BANG
		if s.match("=") {
			b = tok.BANG_EQUAL
		}
		s.addToken(b)
		break
	case "=":
		e := tok.EQUAL
		if s.match("=") {
			e = tok.EQUAL_EQUAL
		}
		s.addToken(e)
		break
	case "<":
		l := tok.LESS
		if s.match("=") {
			l = tok.LESS_EQUAL
		}
		s.addToken(l)
		break
	case ">":
		g := tok.GREATER
		if s.match("=") {
			g = tok.GREATER_EQUAL
		}
		s.addToken(g)
		break
	case "/":
		if s.match("/") {
			for s.peek() != "\n" && !s.isAtEnd() {
				s.advance()
			}
		} else {
			s.addToken(tok.SLASH)
		}
		break
	case " ":
	case "\r":
	case "\t":
		// Ignore
		break
	case "\n":
		s.line++
		break
	case "\"":
		s.string()
		break
	default:
		if s.isDigit(c) {
			s.number()
		} else if s.isAlpha(c) {
			s.identifier()
		} else {
			s.Error(s.line, "Unexpected character.")
		}
		break
	}
}

func (s *Scanner) addToken(tokenType tok.TokenType, literals ...interface{}) {
	text := s.source[s.start:s.current]
	var literal interface{}
	if literals != nil {
		literal = literals[0]
	}
	s.tokens = append(s.tokens, tok.Token{tokenType, text, literal, s.line})
}

func (s *Scanner) advance() byte {
	s.current++
	return s.source[s.current-1]
}

func (s *Scanner) match(expected string) bool {
	if s.isAtEnd() {
		return false
	}
	if s.source[s.current] != expected[0] {
		return false
	}
	s.current++
	return true
}

func (s *Scanner) peek() string {
	if s.isAtEnd() {
		return "\000"
	}
	return string(s.source[s.current])
}

func (s *Scanner) peekNext() string {
	if s.current+1 >= len(s.source) {
		return "\000"
	}
	return string(s.source[s.current+1])
}

func (s *Scanner) string() {
	for s.peek() != "\"" && !s.isAtEnd() {
		if s.peek() == "\n" {
			s.line++
		}
		s.advance()
	}

	// We never reached another "
	if s.isAtEnd() {
		s.Error(s.line, "Unterminated string.")
		return
	}

	// We are at closing "
	s.advance()

	// Get the whole string at once
	str := s.source[s.start+1 : s.current-1]
	s.addToken(tok.STRING, str)
}

func (s *Scanner) number() {
	for s.isDigit(s.peek()) {
		s.advance()
	}

	// . ?
	if s.peek() == "." && s.isDigit(s.peekNext()) {
		s.advance()
		for s.isDigit(s.peek()) {
			s.advance()
		}
	}

	num, err := strconv.ParseFloat(s.source[s.start:s.current], 64)
	if err != nil {
		s.Error(s.line, "Could not convert to number.")
		return
	}
	s.addToken(tok.NUMBER, num)
}

func (s *Scanner) identifier() {
	for s.isAlphaNumeric(s.peek()) {
		s.advance()
	}
	txt := s.source[s.start:s.current]
	tt, found := keywords[txt]
	if !found {
		tt = tok.IDENTIFIER
	}
	s.addToken(tt)
}

func (s Scanner) isDigit(c string) bool {
	return c >= "0" && c <= "9"
}

func (s Scanner) isAlpha(c string) bool {
	return (c >= "a" && c <= "z") ||
		(c >= "A" && c <= "Z") ||
		c == "_"
}

func (s Scanner) isAlphaNumeric(c string) bool {
	return s.isAlpha(c) || s.isDigit(c)
}

func (s Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

// Error spits out failures in the program
func (s *Scanner) Error(line int, message string) {
	s.report(line, "", message)
}

func (s *Scanner) report(line int, where, message string) {
	fmt.Printf("[line \"%d\"] Error %s: %s\n", line, where, message)
	s.HadError = true
}
