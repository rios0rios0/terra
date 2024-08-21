package entities

import "github.com/brianvoe/gofakeit/v7"

type DefaultBaseBuilder[T any] struct {
	modifiers []func(*T)
}

func (it *DefaultBaseBuilder[T]) Build() T {
	var entity T
	_ = gofakeit.Struct(&entity)
	for _, modifier := range it.modifiers {
		modifier(&entity)
	}
	return entity
}

func (it *DefaultBaseBuilder[T]) BuildMany() []T {
	return it.BuildManyDefined(gofakeit.Number(1, 5))
}

func (it *DefaultBaseBuilder[T]) BuildManyDefined(quantity int) []T {
	entities := make([]T, quantity)
	for i := 0; i < quantity; i++ {
		entities[i] = it.Build()
	}
	return entities
}

func (it *DefaultBaseBuilder[T]) AppendModifier(modifier func(*T)) {
	it.modifiers = append(it.modifiers, modifier)
}
