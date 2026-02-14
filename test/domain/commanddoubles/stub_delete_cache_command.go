//go:build integration || unit || test

package commanddoubles //nolint:staticcheck // Test package naming follows established project structure

// StubDeleteCacheCommand is a stub implementation of the DeleteCache interface.
type StubDeleteCacheCommand struct {
	ExecuteCallCount int
	LastToBeDeleted  []string
	LastGlobal       bool
}

func (m *StubDeleteCacheCommand) Execute(toBeDeleted []string, global bool) {
	m.ExecuteCallCount++
	m.LastToBeDeleted = toBeDeleted
	m.LastGlobal = global
}
