package commands

type RunFromRoot interface {
	Execute(toBeDeleted []string)
}
