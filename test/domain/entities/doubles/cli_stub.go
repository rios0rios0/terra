package doubles

import "github.com/brianvoe/gofakeit/v7"

type CLIStub struct {
	canChange bool
}

func NewCLIStub() *CLIStub {
	return &CLIStub{}
}

func (it *CLIStub) WithCanChangeAccount(canChange bool) *CLIStub {
	it.canChange = canChange
	return it
}

func (it *CLIStub) GetName() string {
	return gofakeit.Name()
}

func (it *CLIStub) CanChangeAccount() bool {
	return it.canChange
}

func (it *CLIStub) GetCommandChangeAccount() []string {
	return []string{"command1", "command2"}
}
