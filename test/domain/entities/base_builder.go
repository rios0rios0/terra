package entities

type BaseBuilder[T any] interface {
	Build() T
	BuildMany() []T
	BuildManyDefined(quantity int) []T
	AppendModifier(modifier func(*T))
}
