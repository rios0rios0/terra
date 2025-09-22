package command_doubles //nolint:revive,staticcheck // Test package naming follows established project structure

// StubDeleteCacheCommand is a stub implementation of the DeleteCache interface.
type StubDeleteCacheCommand struct {
	ExecuteCallCount int
	LastToBeDeleted  []string
}

func (m *StubDeleteCacheCommand) Execute(toBeDeleted []string) {
	m.ExecuteCallCount++
	m.LastToBeDeleted = toBeDeleted
}
