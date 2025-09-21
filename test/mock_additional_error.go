package test

// MockAdditionalError implements the error interface.
type MockAdditionalError struct {
	message string
}

func (e *MockAdditionalError) Error() string {
	return e.message
}