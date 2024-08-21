package doubles

type DeleteCacheCommandStub struct {
	success bool
}

func NewDeleteCacheCommandStub() *DeleteCacheCommandStub {
	return &DeleteCacheCommandStub{}
}

func (it *DeleteCacheCommandStub) WithSuccess() *DeleteCacheCommandStub {
	it.success = true
	return it
}

func (it *DeleteCacheCommandStub) WithError() *DeleteCacheCommandStub {
	it.success = false
	return it
}

func (it *DeleteCacheCommandStub) Execute(_ []string) {
	it.success = true
}
