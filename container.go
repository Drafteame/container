package container

import (
	"github.com/Drafteame/container/dependency"
	"github.com/Drafteame/container/injector"
)

var depContainer Container

// Container represents a dependency container that should register factory methods and its dependency threes to be
// injected when
type Container interface {
	Override(name string, dep dependency.Dependency) error
	Provide(name string, dep dependency.Dependency) error
	Get(name string) (any, error)
	Flush()
	Remove(name string)
}

// get returns a global instance for the dependency injection container. If the container is nil, then it will initialize
// a new instance before returning the container.
func get() Container {
	if depContainer == nil {
		depContainer = injector.New()
	}

	return depContainer
}

// New Return a new isolated instance for the dependency injection container. This instance is totally different from
// the global container and does not share any saved dependency three between each other.
func New() Container {
	return injector.New()
}

// Flush WARNING: This function will delete all saved instances, solved and registered factories from the container.
// Do not use this method on production and just use it for testing purposes.
func Flush() {
	get().Flush()
}

// Remove WARNING: This function will remove a specific factory and its solved dependency from the container. Do not use
// this method on production and just use it for testing purposes.
func Remove(name string) {
	get().Remove(name)
}
