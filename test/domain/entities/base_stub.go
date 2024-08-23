package entities

import (
	"encoding/json"
	"fmt"
)

type BaseStub[T any] struct {
	Err error
}

func (it *BaseStub[T]) WithSuccess() *T {
	it.Err = nil
	return it.asT()
}

func (it *BaseStub[T]) WithError() *T {
	return it.WithErrorDefined(fmt.Errorf("any error"))
}

func (it *BaseStub[T]) WithErrorDefined(err error) *T {
	it.Err = err
	return it.asT()
}

// hacky way to convert the struct to the desired type
func (it *BaseStub[T]) asT() *T {
	var result T
	itself, _ := json.Marshal(&it)
	_ = json.Unmarshal(itself, &result)
	return &result
}
