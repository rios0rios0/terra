package command_doubles

// StubSelfUpdateCommand provides a stub implementation for SelfUpdate interface for testing.
type StubSelfUpdateCommand struct {
	ExecuteError   error
	ExecuteCalled  bool
	DryRunFlag     bool
	ForceFlag      bool
	CallCount      int
}

// NewStubSelfUpdateCommand creates a new stub self-update command.
func NewStubSelfUpdateCommand() *StubSelfUpdateCommand {
	return &StubSelfUpdateCommand{}
}

// Execute implements the SelfUpdate interface.
func (s *StubSelfUpdateCommand) Execute(dryRun, force bool) error {
	s.ExecuteCalled = true
	s.CallCount++
	s.DryRunFlag = dryRun
	s.ForceFlag = force
	return s.ExecuteError
}

// WasCalled returns true if Execute was called.
func (s *StubSelfUpdateCommand) WasCalled() bool {
	return s.ExecuteCalled
}

// Reset resets the stub to initial state.
func (s *StubSelfUpdateCommand) Reset() {
	s.ExecuteCalled = false
	s.DryRunFlag = false
	s.ForceFlag = false
	s.CallCount = 0
	s.ExecuteError = nil
}