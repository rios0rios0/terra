package doubles

import "fmt"

type ShellRepositoryStub struct {
	error error
}

func NewShellRepositoryStub() *ShellRepositoryStub {
	return &ShellRepositoryStub{}
}

func (it *ShellRepositoryStub) WithSuccess() *ShellRepositoryStub {
	it.error = nil
	return it
}

func (it *ShellRepositoryStub) WithError() *ShellRepositoryStub {
	it.error = fmt.Errorf("failed to perform command execution")
	return it
}

func (it *ShellRepositoryStub) ExecuteCommand(_ string, _ []string, _ string) error {
	return it.error
}
