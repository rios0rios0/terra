package command_doubles

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
