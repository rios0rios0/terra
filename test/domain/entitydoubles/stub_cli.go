//go:build integration || unit || test

package entitydoubles //nolint:staticcheck // Test package naming follows established project structure

// StubCLI is a stub implementation of entities.CLI.
type StubCLI struct {
	Name                  string
	CanChangeAccountValue bool
	CommandChangeAccount  []string
}

func (m *StubCLI) GetName() string {
	return m.Name
}

func (m *StubCLI) CanChangeAccount() bool {
	return m.CanChangeAccountValue
}

func (m *StubCLI) GetCommandChangeAccount() []string {
	return m.CommandChangeAccount
}
