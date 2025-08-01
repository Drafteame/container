package container

import (
	"errors"

	"github.com/Drafteame/container/dependency"
)

var (
	ErrFactoryNotFunction = errors.New("factory parameter should be a function or a dependency.Dependency instance")
)

// Register It adds a new injection dependency to the container, getting the first result type of the constructor to
// associate the constructor on the injection dependency threes.
//
// This injection will be resolved and built on execution time when the `inject.Invoke(...)` or `inject.Get(name)`
// methods are called.
func Register(name string, factory any, opts ...Option) error {
	return registerDep(name, false, factory, opts...)
}

// MustRegister registers a dependency with the given name and factory. Panics if registration fails.
func MustRegister(name string, factory any, opts ...Option) {
	err := registerDep(name, false, factory, opts...)
	if err != nil {
		panic(err)
	}
}

// Singleton It adds a new injection dependency to the container, getting the first result type of the constructor to
// associate the constructor on the injection dependency threes as a singleton instance.
//
// This function also receives dependency arguments as variadic in case the factory was a function instead of a
// dependency.Dependency.
func Singleton(name string, factory any, opts ...Option) error {
	return registerDep(name, true, factory, opts...)
}

// MustSingleton registers a dependency with the given name and factory as a singleton. Panics if registration fails.
func MustSingleton(name string, factory any, opts ...Option) {
	err := registerDep(name, true, factory, opts...)
	if err != nil {
		panic(err)
	}
}

// Override Set a new dependency that replaces the old one to change behavior on runtime.
// WARNING: This function will remove a specific factory and its solved dependency from the container. Do not use
// this method on production and use it for testing purposes.
func Override(name string, factory any, opts ...Option) error {
	depOpts := buildOptions(opts...)
	c := depOpts.container

	return c.Override(name, dependency.New(factory, depOpts.args...))
}

// MustOverride Set a new dependency that replaces the old one to change behavior on runtime.
// WARNING: This function will remove a specific factory and its solved dependency from the container. Do not use
// this method on production and use it for testing purposes.
func MustOverride(name string, factory any, opts ...Option) {
	err := Override(name, factory, opts...)
	if err != nil {
		panic(err)
	}
}

// Inject is a Wrapper over the dependency.Inject function to generify string symbol name.
func Inject(name string) dependency.Injectable {
	return dependency.Inject(name)
}

func registerDep(name string, singleton bool, factory any, opts ...Option) error {
	depOpts := buildOptions(opts...)

	c := depOpts.container

	if dep, ok := factory.(dependency.Dependency); ok {
		dep.Singleton = singleton
		return c.Provide(name, dep)
	}

	if _, ok := factory.(dependency.Builder); ok {
		return ErrFactoryNotFunction
	}

	if singleton {
		return c.Provide(name, dependency.NewSingleton(factory, depOpts.args...))
	}

	return c.Provide(name, dependency.New(factory, depOpts.args...))
}
