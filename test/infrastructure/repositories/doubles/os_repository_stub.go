package doubles

import (
	"github.com/rios0rios0/terra/internal/domain/entities"
	testentities "github.com/rios0rios0/terra/test/domain/entities"
)

type OSRepositoryStub struct {
	testentities.BaseStub[OSRepositoryStub]
	userID int
}

func NewOSRepositoryStub() *OSRepositoryStub {
	return &OSRepositoryStub{}
}

func (it *OSRepositoryStub) WithUserID(userID int) *OSRepositoryStub {
	it.userID = userID
	return it
}

func (it *OSRepositoryStub) IsSuperUser() bool {
	return it.userID == 0
}

func (it *OSRepositoryStub) ExecuteCommand(_ string, _ []string, _ string) error {
	return it.Err
}

func (it *OSRepositoryStub) InstallExecutable(_, _ string, _ entities.OS) error {
	return it.Err
}
