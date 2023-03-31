package injector

import (
	"github.com/Drafteame/container/dependency"
	"github.com/Drafteame/container/types"
)

// Container is a dependency injection Container implementation
type Container struct {
	solvedDeps map[types.Symbol]any
	deps       map[types.Symbol]dependency.Dependency
}

// New creates a new instance of a Container.
func New() *Container {
	return &Container{
		solvedDeps: make(map[types.Symbol]any),
		deps:       make(map[types.Symbol]dependency.Dependency),
	}
}

// Flush WARNING: This function will delete all saved instances, solved and registered factories from the container.
// Do not use this method on production, and just use it on testing purposes.
func (c *Container) Flush() {
	c.solvedDeps = make(map[types.Symbol]any)
	c.deps = make(map[types.Symbol]dependency.Dependency)
}

// Remove WARNING: This function will remove a specific factory and its solve dependency from the container. Do not use
// this method on production, and just use it on testing purposes.
func (c *Container) Remove(name types.Symbol) {
	delete(c.solvedDeps, name)
	delete(c.deps, name)
}
