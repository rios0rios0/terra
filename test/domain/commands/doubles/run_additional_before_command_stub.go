package doubles

import "github.com/rios0rios0/terra/internal/domain/commands/interfaces"

type RunAdditionalBeforeCommandStub struct {
	success bool
}

func NewRunAdditionalBeforeCommandStub() *RunAdditionalBeforeCommandStub {
	return &RunAdditionalBeforeCommandStub{}
}

func (it *RunAdditionalBeforeCommandStub) WithSuccess() *RunAdditionalBeforeCommandStub {
	it.success = true
	return it
}

func (it *RunAdditionalBeforeCommandStub) WithError() *RunAdditionalBeforeCommandStub {
	it.success = false
	return it
}

func (it *RunAdditionalBeforeCommandStub) Execute(_ string, _ []string, _ interfaces.RunAdditionalBeforeListeners) {
	it.success = true
}
