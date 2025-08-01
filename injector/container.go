package injector

import (
	"github.com/Drafteame/container/dependency"
)

// Container is a dependency injection Container implementation
type Container struct {
	solvedDeps map[string]any
	deps       map[string]dependency.Dependency
}

// New creates a new instance of a Container.
func New() *Container {
	return &Container{
		solvedDeps: make(map[string]any),
		deps:       make(map[string]dependency.Dependency),
	}
}

// Flush WARNING: This function will delete all saved instances, solved and registered factories from the container.
// Do not use this method on production and just use it for testing purposes.
func (c *Container) Flush() {
	c.solvedDeps = make(map[string]any)
	c.deps = make(map[string]dependency.Dependency)
}

// Remove WARNING: This function will remove a specific factory and its solved dependency from the container. Do not use
// this method on production and use it for testing purposes.
func (c *Container) Remove(name string) {
	delete(c.solvedDeps, name)
	delete(c.deps, name)
}

// Override Set a new dependency that replaces the old one to change behavior on runtime.
// WARNING: This function will remove a specific factory and its solved dependency from the container. Do not use
// this method on production and use it for testing purposes.
func (c *Container) Override(name string, dep dependency.Dependency) error {
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
