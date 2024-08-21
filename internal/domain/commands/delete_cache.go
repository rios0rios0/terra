package commands

type DeleteCache interface {
	Execute(toBeDeleted []string)
}
