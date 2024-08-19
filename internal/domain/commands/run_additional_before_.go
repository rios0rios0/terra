package commands

type RunAdditionalBefore interface {
	Execute(targetPath string, arguments []string)
}
