package errors

import "fmt"

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
