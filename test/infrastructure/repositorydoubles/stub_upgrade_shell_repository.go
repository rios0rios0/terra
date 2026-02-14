//go:build integration || unit || test

package repositorydoubles //nolint:staticcheck // Test package naming follows established project structure

// StubUpgradeShellRepository is a stub implementation of the UpgradeShellRepository interface.
type StubUpgradeShellRepository struct {
	ExecuteCallCount int
	LastCommand      string
	LastArguments    []string
	LastDirectory    string
	ErrorToReturn    error
}

func (m *StubUpgradeShellRepository) ExecuteCommandWithUpgrade(
	command string,
	arguments []string,
	directory string,
) error {
	m.ExecuteCallCount++
	m.LastCommand = command
	m.LastArguments = arguments
	m.LastDirectory = directory

	return m.ErrorToReturn
}
