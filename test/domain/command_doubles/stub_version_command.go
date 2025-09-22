package command_doubles //nolint:revive,staticcheck // Test package naming follows established project structure

// StubVersionCommand is a stub implementation of the Version interface.
type StubVersionCommand struct {
	ExecuteCallCount int
}

func (m *StubVersionCommand) Execute() {
	m.ExecuteCallCount++
}
