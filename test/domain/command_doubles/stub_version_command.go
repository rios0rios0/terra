//go:build integration || unit || test

package command_doubles //nolint:staticcheck // Test package naming follows established project structure

// StubVersionCommand is a stub implementation of the Version interface.
type StubVersionCommand struct {
	ExecuteCallCount int
}

func (m *StubVersionCommand) Execute() {
	m.ExecuteCallCount++
}
