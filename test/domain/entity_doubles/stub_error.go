package entity_doubles

// StubError implements the error interface for testing
type StubError struct {
	message string
}

func NewStubError(message string) *StubError {
	return &StubError{message: message}
}

func (e *StubError) Error() string {
	return e.message
}