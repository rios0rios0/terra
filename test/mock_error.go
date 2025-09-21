package test

// MockError implements the error interface for testing
type MockError struct {
	message string
}

func (e *MockError) Error() string {
	return e.message
}