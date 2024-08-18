package commands

type ClearCache interface {
	Execute(toBeDeleted []string)
}
