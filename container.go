package container

import (
	"github.com/Drafteame/container/dependency"
	"github.com/Drafteame/container/injector"
	"github.com/Drafteame/container/types"
)

var depContainer Container

type symbolName interface {
	string | types.Symbol
}

// Container represents a dependency container that should register factory methods and its dependency threes to be
// injected when
type Container interface {
	Provide(name types.Symbol, dep dependency.Dependency) error
	Invoke(construct any) error
	Get(name types.Symbol) (any, error)
	Flush()
	Remove(name types.Symbol)
}

// get return a global instance for the dependency injection container. If the container is nil, then it will initialize
// a new instance before returning the container.
func get() Container {
	if depContainer == nil {
		depContainer = injector.New()
	}

	return depContainer
}

// New Return a new isolated instance for the dependency injection container. This instance is totally different from
// the global container and do not share any saved dependency three between each other.
func New() Container {
	return injector.New()
}

// Invoke Is the entry point to execute dependency injection resolution. It calls an invoker function that can
// receive or not a struct that embeds inject.In struct as input, and return an error or not (any other return field or
// type will be ignored on resolution). When invoker is called it will resolve the dependency threes of each field from
// the previously provided resources on Container.
func Invoke(construct any) error {
	return get().Invoke(construct)
}

// Flush WARNING: This function will delete all saved instances, solved and registered factories from the container.
// Do not use this method on production, and just use it for testing purposes.
func Flush() {
	get().Flush()
}

// Remove WARNING: This function will remove a specific factory and its solve dependency from the container. Do not use
// this method on production, and just use it on testing purposes.
func Remove[T symbolName](name T) {
	get().Remove(types.Symbol(name))
}
