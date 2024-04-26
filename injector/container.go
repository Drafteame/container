package injector

import (
	"github.com/Drafteame/container/dependency"
	"github.com/Drafteame/container/types"
)

// Container is a dependency injection Container implementation
type Container struct {
	solvedDeps map[types.Symbol]any
	deps       map[types.Symbol]dependency.Dependency
	testMode   bool
}

// New creates a new instance of a Container.
func New() *Container {
	return &Container{
		solvedDeps: make(map[types.Symbol]any),
		deps:       make(map[types.Symbol]dependency.Dependency),
		testMode:   false,
	}
}

// TestMode WARNING: Sets testMode flag to true, bypassing singleton instance generation to avoid race conditions when
// container is used on test cases.
func (c *Container) TestMode() {
	c.testMode = true
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

// Override Set a new dependency that replaces the old one to change behavior on runtime.
// WARNING: This function will remove a specific factory and its solve dependency from the container. Do not use
// this method on production, and just use it on testing purposes.
func (c *Container) Override(name types.Symbol, dep dependency.Dependency) error {
	previous, ok := c.deps[name]
	if !ok {
		return c.Provide(name, dep)
	}

	if previous.IsSingleton() {
		dep.Singleton = true
	}

	c.Remove(name)
	return c.Provide(name, dep)
}
