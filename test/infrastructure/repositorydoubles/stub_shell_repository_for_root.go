//go:build integration || unit || test

package repositorydoubles //nolint:staticcheck // Test package naming follows established project structure

// StubShellRepositoryForRoot is a stub implementation of repositories.ShellRepositoryWithUpgrade.
type StubShellRepositoryForRoot struct {
	ExecuteCallCount int
	LastCommand      string
	LastArguments    []string
	LastDirectory    string
	ExecuteErrors    []error
	CallHistory      []struct {
		Command   string
		Arguments []string
		Directory string
	}
}

func (m *StubShellRepositoryForRoot) ExecuteCommand(
	command string,
	arguments []string,
	directory string,
) error {
	m.CallHistory = append(m.CallHistory, struct {
		Command   string
		Arguments []string
		Directory string
	}{
		Command:   command,
		Arguments: arguments,
		Directory: directory,
	})

	m.ExecuteCallCount++
	m.LastCommand = command
	m.LastArguments = arguments
	m.LastDirectory = directory

	if len(m.ExecuteErrors) > 0 {
		err := m.ExecuteErrors[0]
		m.ExecuteErrors = m.ExecuteErrors[1:]
		return err
	}
	return nil
}

// ExecuteCommandWithUpgradeDetection implements the ShellRepositoryWithUpgrade interface
func (m *StubShellRepositoryForRoot) ExecuteCommandWithUpgradeDetection(
	command string,
	arguments []string,
	directory string,
) error {
	// For testing purposes, delegate to ExecuteCommand
	return m.ExecuteCommand(command, arguments, directory)
}
