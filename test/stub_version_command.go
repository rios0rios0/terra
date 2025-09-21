package test

// StubVersionCommand is a stub implementation of the Version interface
type StubVersionCommand struct {
	ExecuteCallCount int
}

func (m *StubVersionCommand) Execute() {
	m.ExecuteCallCount++
}