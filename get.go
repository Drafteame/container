package container

import (
	"fmt"
	"reflect"

	"github.com/Drafteame/container/types"
)

// Get is a wrapper over the Get function attached to the global container. This function modify the return type of the
// resolved dependency, returned as `any` to the provided generic type `T`. If it can't be casted it will return an
// error.
func Get[T any, K symbolName](name K) (T, error) {
	instance, err := get().Get(types.Symbol(name))
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

// MustGet Same functionality that Get function but instead of returning error, it panics.
func MustGet[T any, K symbolName](name K) T {
	instance, err := get().Get(types.Symbol(name))
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
