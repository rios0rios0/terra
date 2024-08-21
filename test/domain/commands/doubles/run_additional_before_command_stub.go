package doubles

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

func (it *RunAdditionalBeforeCommandStub) Execute(_ string, _ []string) {
	it.success = true
}
