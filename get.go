package container

import (
	"fmt"
	"reflect"
)

// Get is a wrapper over the Get function attached to the global container. This function modifies the return type of the
// resolved dependency, returned as `any` to the provided generic type `T`. If it can't be cast it will return an
// error.
func Get[T any](name string, opts ...Option) (T, error) {
	depOpts := buildOptions(opts...)

	c := depOpts.container

	instance, err := c.Get(name)
	if err != nil {
		aux := new(T)
		return *aux, err
	}

	cast, ok := instance.(T)
	if !ok {
		aux := new(T)
		axtype := reflect.TypeOf(*aux)
		return *aux, fmt.Errorf("inject: error casting instance of `%s` dependency to `%v`", name, axtype)
	}

	return cast, nil
}

// MustGet Same functionality that Get function, but instead of returning an error, it panics.
func MustGet[T any](name string, opts ...Option) T {
	depOpts := buildOptions(opts...)

	c := depOpts.container

	instance, err := c.Get(name)
	if err != nil {
		panic(err)
	}

	cast, ok := instance.(T)
	if !ok {
		aux := new(T)
		axtype := reflect.TypeOf(*aux)

		errCast := fmt.Errorf("inject: error casting instance of `%s` dependency to `%v`", name, axtype)
		panic(errCast)
	}

	return cast
}
