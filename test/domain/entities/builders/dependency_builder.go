package builders

import (
	"github.com/rios0rios0/terra/internal/domain/entities"
	testentities "github.com/rios0rios0/terra/test/domain/entities"
)

type DependencyBuilder struct {
	testentities.BaseBuilder[entities.Dependency]
}

func NewDependencyBuilder() *DependencyBuilder {
	return &DependencyBuilder{}
}

func (it *DependencyBuilder) WithName(name string) *DependencyBuilder {
	it.AppendModifier(func(entity *entities.Dependency) {
		entity.Name = name
	})
	return it
}

func (it *DependencyBuilder) WithCLI(cli string) *DependencyBuilder {
	it.AppendModifier(func(entity *entities.Dependency) {
		entity.CLI = cli
	})
	return it
}

func (it *DependencyBuilder) WithVersionURL(versionURL string) *DependencyBuilder {
	it.AppendModifier(func(entity *entities.Dependency) {
		entity.VersionURL = versionURL
	})
	return it
}

func (it *DependencyBuilder) WithBinaryURL(binaryURL string) *DependencyBuilder {
	it.AppendModifier(func(entity *entities.Dependency) {
		entity.BinaryURL = binaryURL
	})
	return it
}

func (it *DependencyBuilder) WithRegexVersion(regexVersion string) *DependencyBuilder {
	it.AppendModifier(func(entity *entities.Dependency) {
		entity.RegexVersion = regexVersion
	})
	return it
}

func (it *DependencyBuilder) WithFormattingCommand(formattingCommand []string) *DependencyBuilder {
	it.AppendModifier(func(entity *entities.Dependency) {
		entity.FormattingCommand = formattingCommand
	})
	return it
}
