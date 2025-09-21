package test

// MockRunAdditionalBefore is a mock implementation for RunAdditionalBefore interface.
type MockRunAdditionalBefore struct {
	ExecuteCalled  bool
	LastTargetPath string
	LastArguments  []string
}

func (m *MockRunAdditionalBefore) Execute(targetPath string, arguments []string) {
	m.ExecuteCalled = true
	m.LastTargetPath = targetPath
	m.LastArguments = arguments
}