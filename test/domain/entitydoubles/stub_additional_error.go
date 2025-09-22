//go:build integration || unit || test

package entitydoubles //nolint:staticcheck // Test package naming follows established project structure

// StubAdditionalError implements the error interface.
type StubAdditionalError struct {
	message string
}

func (e *StubAdditionalError) Error() string {
	return e.message
}
