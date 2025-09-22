package entity_doubles //nolint:revive,staticcheck // Test package naming follows established project structure

// StubRootError implements the error interface.
type StubRootError struct {
	message string
}

func (e *StubRootError) Error() string {
	return e.message
}
