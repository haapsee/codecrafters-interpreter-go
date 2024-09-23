package errors

type LexicalError struct {
	Message string
}

func (le LexicalError) Error() string {
	return le.Message
}

func NewLexicalError(line int, where string, message string) error {
	return LexicalError{
		Message: message,
	}
}
