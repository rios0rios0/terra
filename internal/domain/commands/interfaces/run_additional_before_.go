package interfaces

type RunAdditionalBefore interface {
	Execute(targetPath string, arguments []string, listeners RunAdditionalBeforeListeners)
}

type RunAdditionalBeforeListeners struct {
	OnSuccess func()
	OnError   func(error)
}
