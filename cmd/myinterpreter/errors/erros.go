package errors

import (
	"fmt"

	"github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/token"
)

type Error struct {
	message string
}

func (e Error) Error() string {
	return e.message
}

func New(message string) Error {
	return Error{
		message: message,
	}
}

type RuntimeError struct {
	Token   token.Token
	Message string
}

func (err RuntimeError) Error() string {
	return fmt.Sprintf("%s\n[line %d]", err.Message, err.Token.Line)
}

func NewRuntimeError(token token.Token, message string) RuntimeError {
	return RuntimeError{
		Message: message,
		Token:   token,
	}
}

type LexicalError struct {
	Message string
	Line    int
	Where   string
}

func (le LexicalError) Error() string {
	return fmt.Sprintf("[line %d] Error: %s", le.Line, le.Message)
}

func NewLexicalError(line int, where string, message string) error {
	return LexicalError{
		Message: message,
		Line:    line,
		Where:   where,
	}
}
