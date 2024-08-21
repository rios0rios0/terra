package doubles

import (
	"github.com/rios0rios0/terra/internal/domain/commands/interfaces"
	"github.com/rios0rios0/terra/internal/domain/entities"
)

type RunFromRootCommandStub struct {
	success bool
}

func NewRunFromRootCommandStub() *RunFromRootCommandStub {
	return &RunFromRootCommandStub{}
}

func (it *RunFromRootCommandStub) WithSuccess() *RunFromRootCommandStub {
	it.success = true
	return it
}

func (it *RunFromRootCommandStub) WithError() *RunFromRootCommandStub {
	it.success = false
	return it
}

func (it *RunFromRootCommandStub) Execute(_ string, _ []string, _ []entities.Dependency, _ interfaces.RunFromRootListeners) {
	it.success = true
}
