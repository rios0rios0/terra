package doubles

import (
	"fmt"
	"github.com/rios0rios0/terra/internal/domain/entities"
)

type OSRepositoryStub struct {
	error error
}

func NewOSRepositoryStub() *OSRepositoryStub {
	return &OSRepositoryStub{}
}

func (it *OSRepositoryStub) WithSuccess() *OSRepositoryStub {
	it.error = nil
	return it
}

func (it *OSRepositoryStub) WithError() *OSRepositoryStub {
	it.error = fmt.Errorf("failed to perform command execution")
	return it
}

func (it *OSRepositoryStub) ExecuteCommand(_ string, _ []string, _ string) error {
	return it.error
}

func (it *OSRepositoryStub) InstallExecutable(_, _ string, _ entities.OS) error {
	return it.error
}
