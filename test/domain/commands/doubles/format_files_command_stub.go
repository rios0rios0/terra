package doubles

import (
	"github.com/rios0rios0/terra/internal/domain/commands/interfaces"
	"github.com/rios0rios0/terra/internal/domain/entities"
)

type FormatFilesCommandStub struct {
	success bool
}

func NewFormatFilesCommandStub() *FormatFilesCommandStub {
	return &FormatFilesCommandStub{}
}

func (it *FormatFilesCommandStub) WithSuccess() *FormatFilesCommandStub {
	it.success = true
	return it
}

func (it *FormatFilesCommandStub) WithError() *FormatFilesCommandStub {
	it.success = false
	return it
}

func (it *FormatFilesCommandStub) Execute(_ []entities.Dependency, _ interfaces.FormatFilesListeners) {
	it.success = true
}
