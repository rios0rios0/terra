package doubles

import "github.com/rios0rios0/terra/internal/domain/entities"

type InstallDependenciesCommandStub struct {
	success bool
}

func NewInstallDependenciesCommandStub() *InstallDependenciesCommandStub {
	return &InstallDependenciesCommandStub{}
}

func (it *InstallDependenciesCommandStub) WithSuccess() *InstallDependenciesCommandStub {
	it.success = true
	return it
}

func (it *InstallDependenciesCommandStub) WithError() *InstallDependenciesCommandStub {
	it.success = false
	return it
}

func (it *InstallDependenciesCommandStub) Execute(_ []entities.Dependency) {
	it.success = true
}
