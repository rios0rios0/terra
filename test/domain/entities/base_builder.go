package entities

import "github.com/brianvoe/gofakeit/v7"

type BaseBuilder[T any] struct {
	modifiers []func(*T)
}

func (it *BaseBuilder[T]) Build() T {
	var entity T
	_ = gofakeit.Struct(&entity)
	for _, modifier := range it.modifiers {
		modifier(&entity)
	}
	return entity
}

func (it *BaseBuilder[T]) BuildMany() []T {
	return it.BuildManyDefined(gofakeit.Number(1, 5))
}

func (it *BaseBuilder[T]) BuildManyDefined(quantity int) []T {
	entities := make([]T, quantity)
	for i := 0; i < quantity; i++ {
		entities[i] = it.Build()
	}
	return entities
}

func (it *BaseBuilder[T]) AppendModifier(modifier func(*T)) {
	it.modifiers = append(it.modifiers, modifier)
}
