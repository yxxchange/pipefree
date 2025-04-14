package interfaces

type BaseUnit[T any] interface {
	Value() T
}

type CompareUnit[T any] interface {
	Less(other CompareUnit[T]) bool
	Equal(other CompareUnit[T]) bool

	BaseUnit[T]
}

type CalculateUnit[T any] interface {
	Add(other CompareUnit[T]) T // +
	Sub(other CompareUnit[T]) T // -
	Div(other CompareUnit[T]) T // /
	Mul(other CompareUnit[T]) T // *
	Mod(other CompareUnit[T]) T // %
}
