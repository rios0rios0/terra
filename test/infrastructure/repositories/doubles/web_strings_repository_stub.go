package doubles

import (
	testentities "github.com/rios0rios0/terra/test/domain/entities"
)

type WebStringsRepositoryStub struct {
	testentities.BaseStub[WebStringsRepositoryStub]
}

func NewWebStringsRepositoryStub() *WebStringsRepositoryStub {
	return &WebStringsRepositoryStub{}
}

func (it *WebStringsRepositoryStub) FindStringMatchInURL(_, _ string) (string, error) {
	return "", it.Err
}
