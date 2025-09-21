package test

// MockCLI is a mock implementation of entities.CLI.
type MockCLI struct {
	Name                  string
	CanChangeAccountValue bool
	CommandChangeAccount  []string
}

func (m *MockCLI) GetName() string {
	return m.Name
}

func (m *MockCLI) CanChangeAccount() bool {
	return m.CanChangeAccountValue
}

func (m *MockCLI) GetCommandChangeAccount() []string {
	return m.CommandChangeAccount
}