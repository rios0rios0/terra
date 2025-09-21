package test

// StubError implements the error interface for testing
type StubError struct {
	message string
}

func (e *StubError) Error() string {
	return e.message
}