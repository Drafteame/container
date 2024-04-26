package container

import (
	"errors"

	"github.com/Drafteame/container/dependency"
	"github.com/Drafteame/container/types"
)

var (
	ErrFactoryNotFunction = errors.New("factory parameter should be a function or a dependency.Dependency instance")
)

// Register It adds a new injection dependency to the container, getting the first result type of the constructor to
// associate the constructor on the injection dependency threes.
//
// This injection will be resolved and built on execution time when the `inject.Invoke(...)` or `inject.Get(name)`
// methods are called.
func Register[T symbolName](name T, factory any, args ...any) error {
	return registerDep(types.Symbol(name), false, factory, args...)
}

// Singleton It adds a new injection dependency to the container, getting the first result type of the constructor to
// associate the constructor on the injection dependency threes as a singleton instance.
//
// This function also receive dependency arguments as variadic in case the factory were a function instead of a
// dependency.Dependency.
func Singleton[T symbolName](name T, factory any, args ...any) error {
	return registerDep(types.Symbol(name), true, factory, args...)
}

// Override Set a new dependency that replaces the old one to change behavior on runtime.
// WARNING: This function will remove a specific factory and its solve dependency from the container. Do not use
// this method on production, and just use it on testing purposes.
func Override[T symbolName](name T, factory any, args ...any) error {
	return get().Override(types.Symbol(name), dependency.New(factory, args...))
}

// Inject is a Wrapper ver the dependency.Inject function to generify string symbol name.
func Inject[T symbolName](name T) dependency.Injectable {
	return dependency.Inject(types.Symbol(name))
}

func registerDep(name types.Symbol, singleton bool, factory any, args ...any) error {
	if dep, ok := factory.(dependency.Dependency); ok {
		dep.Singleton = singleton
		return get().Provide(name, dep)
	}

	if _, ok := factory.(dependency.Builder); ok {
		return ErrFactoryNotFunction
	}

	if singleton {
		return get().Provide(name, dependency.NewSingleton(factory, args...))
	}

	return get().Provide(name, dependency.New(factory, args...))
}
