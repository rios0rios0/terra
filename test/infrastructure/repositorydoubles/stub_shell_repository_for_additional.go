package repositorydoubles

// StubShellRepositoryForAdditional is a stub implementation of repositories.ShellRepository.
type StubShellRepositoryForAdditional struct {
	ExecuteCallCount int
	LastCommand      string
	LastArguments    []string
	LastDirectory    string
	ExecuteErrors    []error
	CallHistory      []CallRecord
}

func (m *StubShellRepositoryForAdditional) ExecuteCommand(
	command string,
	arguments []string,
	directory string,
) error {
	m.CallHistory = append(m.CallHistory, CallRecord{
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
