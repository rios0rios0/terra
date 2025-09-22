package entity_doubles

// StubRootError implements the error interface.
type StubRootError struct {
	message string
}

func (e *StubRootError) Error() string {
	return e.message
}