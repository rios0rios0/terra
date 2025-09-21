package test

// StubAdditionalError implements the error interface.
type StubAdditionalError struct {
	message string
}

func (e *StubAdditionalError) Error() string {
	return e.message
}