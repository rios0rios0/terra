package test

// MockDeleteCacheCommand is a mock implementation of the DeleteCache interface
type MockDeleteCacheCommand struct {
	ExecuteCallCount int
	LastToBeDeleted  []string
}

func (m *MockDeleteCacheCommand) Execute(toBeDeleted []string) {
	m.ExecuteCallCount++
	m.LastToBeDeleted = toBeDeleted
}