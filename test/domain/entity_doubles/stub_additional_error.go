//nolint:revive,staticcheck // Test package naming follows established project structure
package entity_doubles

// StubAdditionalError implements the error interface.
type StubAdditionalError struct {
	message string
}

func (e *StubAdditionalError) Error() string {
	return e.message
}
