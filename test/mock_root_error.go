package test

// MockRootError implements the error interface.
type MockRootError struct {
	message string
}

func (e *MockRootError) Error() string {
	return e.message
}