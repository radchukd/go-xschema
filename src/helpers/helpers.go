package helpers

type XObject interface {
	Validate(interface{}) (bool, []error)
	String() string
}

type XValidation[T any] struct {
	E error
	F func(T) bool
}
