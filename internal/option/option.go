package option

// Option is an optional value.
type Option[T any] struct {
	value T
	set   bool
}

// Some returns an actual value wrapped in option.
func Some[T any](value T) Option[T] {
	return Option[T]{
		value: value,
		set:   true,
	}
}

// None returns an empty option.
func None[T any]() Option[T] {
	return Option[T]{}
}

// Value returns the value.
func (o *Option[T]) Value() T {
	if !o.set {
		panic("value not set")
	}
	return o.value
}

// HasValue returns true if the option has a value.
func (o *Option[T]) HasValue() bool {
	return o.set
}
