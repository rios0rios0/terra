package test

// MockVersionCommand is a mock implementation of the Version interface
type MockVersionCommand struct {
	ExecuteCallCount int
}

func (m *MockVersionCommand) Execute() {
	m.ExecuteCallCount++
}