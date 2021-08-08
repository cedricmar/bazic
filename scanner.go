package main

import "strconv"

type Scanner struct {
	source  string
	tokens  []Token
	start   int
	current int
	line    int
}

var keywords = map[string]TokenType{
	"and":    AND,
	"class":  CLASS,
	"else":   ELSE,
	"false":  FALSE,
	"for":    FOR,
	"fun":    FUN,
	"if":     IF,
	"nil":    NIL,
	"or":     OR,
	"print":  PRINT,
	"return": RETURN,
	"super":  SUPER,
	"this":   THIS,
	"true":   TRUE,
	"var":    VAR,
	"while":  WHILE,
}

func NewScanner(source string) Scanner {
	return Scanner{source, []Token{}, 0, 0, 1}
}

func (s *Scanner) ScanTokens() []Token {
	for !s.isAtEnd() {
		s.start = s.current
		s.scanToken()
	}

	s.tokens = append(s.tokens, Token{EOF, "", "", s.line})
	return s.tokens
}

func (s *Scanner) scanToken() {
	c := string(s.advance())
	switch c {
	case "(":
		s.addToken(LEFT_PAREN)
		break
	case ")":
		s.addToken(RIGHT_PAREN)
		break
	case "{":
		s.addToken(LEFT_BRACE)
		break
	case "}":
		s.addToken(RIGHT_BRACE)
		break
	case ",":
		s.addToken(COMMA)
		break
	case ".":
		s.addToken(DOT)
		break
	case "-":
		s.addToken(MINUS)
		break
	case "+":
		s.addToken(PLUS)
		break
	case ";":
		s.addToken(SEMICOLON)
		break
	case "*":
		s.addToken(STAR)
		break
	case "!":
		b := BANG
		if s.match("=") {
			b = BANG_EQUAL
		}
		s.addToken(b)
		break
	case "=":
		e := EQUAL
		if s.match("=") {
			e = EQUAL_EQUAL
		}
		s.addToken(e)
		break
	case "<":
		l := LESS
		if s.match("=") {
			l = LESS_EQUAL
		}
		s.addToken(l)
		break
	case ">":
		g := GREATER
		if s.match("=") {
			g = GREATER_EQUAL
		}
		s.addToken(g)
		break
	case "/":
		if s.match("/") {
			for s.peek() != "\n" && !s.isAtEnd() {
				s.advance()
			}
		} else {
			s.addToken(SLASH)
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
			Error(s.line, "Unexpected character.")
		}
		break
	}
}

func (s *Scanner) addToken(tokenType TokenType, literals ...interface{}) {
	text := s.source[s.start:s.current]
	var literal interface{}
	if literals != nil {
		literal = literals[0]
	}
	s.tokens = append(s.tokens, Token{tokenType, text, literal, s.line})
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
		Error(s.line, "Unterminated string.")
		return
	}

	// We are at closing "
	s.advance()

	// Get the whole string at once
	str := s.source[s.start+1 : s.current-1]
	s.addToken(STRING, str)
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
		Error(s.line, "Could not convert to number.")
		return
	}
	s.addToken(NUMBER, num)
}

func (s *Scanner) identifier() {
	for s.isAlphaNumeric(s.peek()) {
		s.advance()
	}
	txt := s.source[s.start:s.current]
	tt, found := keywords[txt]
	if !found {
		tt = IDENTIFIER
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
