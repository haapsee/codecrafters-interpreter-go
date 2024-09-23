package scanner

import (
	"fmt"
	"strconv"

	"github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/errors"
	"github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/functions"
	"github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/token"
)

type Scanner struct {
	CurrentIndex int
	HadError     bool
	Line         int
	Source       string
	StartIndex   int
	Tokens       []token.Token
}

func (s *Scanner) isAtEnd() bool {
	return s.CurrentIndex >= len(s.Source)
}

func (s *Scanner) getCurrentSubString() string {
	return s.Source[s.StartIndex:s.CurrentIndex]
}

func (s *Scanner) addToken(tokentype token.TokenType, literal interface{}) {
	newToken := token.NewToken(tokentype, s.getCurrentSubString(), literal, s.Line)
	s.Tokens = append(s.Tokens, newToken)
}

func (s *Scanner) peekNext() rune {
	if s.CurrentIndex+1 >= len(s.Source) {
		return rune(0)
	}
	return []rune(s.Source)[s.CurrentIndex+1]
}

func (s *Scanner) peek() rune {
	if s.isAtEnd() {
		return rune(0)
	}
	return []rune(s.Source)[s.CurrentIndex]
}

func (s *Scanner) match(char rune) bool {
	if s.isAtEnd() || s.peek() != char {
		return false
	}
	s.CurrentIndex++
	return true
}

func (s *Scanner) advance() rune {
	char := []rune(s.Source)[s.CurrentIndex]
	s.CurrentIndex++
	return char
}

func (s *Scanner) checkReservedIdentifiers() bool {
	switch s.getCurrentSubString() {
	case "and":
		s.addToken(token.AND, nil)
	case "class":
		s.addToken(token.CLASS, nil)
	case "else":
		s.addToken(token.ELSE, nil)
	case "false":
		s.addToken(token.FALSE, nil)
	case "for":
		s.addToken(token.FOR, nil)
	case "fun":
		s.addToken(token.FUN, nil)
	case "if":
		s.addToken(token.IF, nil)
	case "nil":
		s.addToken(token.NIL, nil)
	case "or":
		s.addToken(token.OR, nil)
	case "print":
		s.addToken(token.PRINT, nil)
	case "return":
		s.addToken(token.RETURN, nil)
	case "super":
		s.addToken(token.SUPER, nil)
	case "this":
		s.addToken(token.THIS, nil)
	case "true":
		s.addToken(token.TRUE, nil)
	case "var":
		s.addToken(token.VAR, nil)
	case "while":
		s.addToken(token.WHILE, nil)
	default:
		return false
	}
	return true
}

func (s *Scanner) parseIdentifier() {
	for ; functions.IsAlphaDigit(s.peek()) && !s.isAtEnd(); s.advance() {
	}

	if !s.checkReservedIdentifiers() {
		s.addToken(token.IDENTIFIER, nil)
	}
}

func (s *Scanner) parseNumber() {
	for ; functions.IsDigit(s.peek()) && !s.isAtEnd(); s.advance() {
	}

	if s.peek() == '.' && functions.IsDigit(s.peekNext()) {
		s.advance()
		for ; functions.IsDigit(s.peek()) && !s.isAtEnd(); s.advance() {
		}
	}

	value, err := strconv.ParseFloat(s.getCurrentSubString(), 64)
	if err != nil {
		panic(err)
	}
	s.addToken(token.NUMBER, value)
}

func (s *Scanner) parseString() error {
	for !s.isAtEnd() && s.peek() != '"' {
		if s.peek() == '\n' {
			s.Line++
		}
		s.advance()
	}

	if s.isAtEnd() {
		return errors.NewLexicalError(s.Line, "", "Unterminated string.")
	}

	result := s.getCurrentSubString()[1:]
	s.advance()
	s.addToken(token.STRING, result)
	return nil
}

func (s *Scanner) scanToken() error {
	char := s.advance()
	switch char {
	case '(':
		s.addToken(token.LEFT_PAREN, nil)
	case ')':
		s.addToken(token.RIGHT_PAREN, nil)
	case '{':
		s.addToken(token.LEFT_BRACE, nil)
	case '}':
		s.addToken(token.RIGHT_BRACE, nil)
	case ',':
		s.addToken(token.COMMA, nil)
	case '.':
		s.addToken(token.DOT, nil)
	case '-':
		s.addToken(token.MINUS, nil)
	case '+':
		s.addToken(token.PLUS, nil)
	case ';':
		s.addToken(token.SEMICOLON, nil)
	case '*':
		s.addToken(token.STAR, nil)
	case '=':
		tokentype := token.EQUAL
		if s.match('=') {
			tokentype = token.EQUAL_EQUAL
		}
		s.addToken(tokentype, nil)
	case '!':
		tokentype := token.BANG
		if s.match('=') {
			tokentype = token.BANG_EQUAL
		}
		s.addToken(tokentype, nil)
	case '<':
		tokentype := token.LESS
		if s.match('=') {
			tokentype = token.LESS_EQUAL
		}
		s.addToken(tokentype, nil)
	case '>':
		tokentype := token.GREATER
		if s.match('=') {
			tokentype = token.GREATER_EQUAL
		}
		s.addToken(tokentype, nil)
	case '/':
		if s.match('/') {
			for ; s.peek() != '\n' && !s.isAtEnd(); s.advance() {
			}
		} else {
			s.addToken(token.SLASH, nil)
		}
	case '"':
		s.parseString()
	case ' ':
	case '\t':
	case '\r':
	case '\n':
		s.Line++
	default:
		if functions.IsDigit(char) {
			s.parseNumber()
		} else if functions.IsAlpha(char) {
			s.parseIdentifier()
		} else {
			message := fmt.Sprintf("Unexpected character: %c", char)
			return errors.NewLexicalError(s.Line, "", message)
		}
	}
	return nil
}

func (s *Scanner) ScanTokens() error {
	for !s.isAtEnd() {
		s.StartIndex = s.CurrentIndex
		err := s.scanToken()
		if err != nil {
			return err
		}
	}
	s.Tokens = append(s.Tokens, token.NewToken(token.EOF, "", nil, s.Line))
	return nil
}

func NewScanner(source string) Scanner {
	return Scanner{
		Source:       source,
		HadError:     false,
		Line:         1,
		StartIndex:   0,
		CurrentIndex: 0,
		Tokens:       []token.Token{},
	}
}
