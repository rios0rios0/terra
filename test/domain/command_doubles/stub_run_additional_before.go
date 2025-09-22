//go:build integration || unit || test

package command_doubles //nolint:staticcheck // Test package naming follows established project structure

// StubRunAdditionalBefore is a stub implementation for RunAdditionalBefore interface.
type StubRunAdditionalBefore struct {
	ExecuteCalled  bool
	LastTargetPath string
	LastArguments  []string
}

func (m *StubRunAdditionalBefore) Execute(targetPath string, arguments []string) {
	m.ExecuteCalled = true
	m.LastTargetPath = targetPath
	m.LastArguments = arguments
}
