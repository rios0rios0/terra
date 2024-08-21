package interfaces

type DeleteCache interface {
	Execute(toBeDeleted []string, listeners DeleteCacheListeners)
}

type DeleteCacheListeners struct {
	OnSuccess func()
	OnError   func(error)
}
