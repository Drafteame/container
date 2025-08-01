package dependency

// Builder definition for a dependency that should be build on injection time.
type Builder interface {
	Build() (any, error)
}

// Container is a container that holds global dependencies.
type Container interface {
	Get(name string) (any, error)
}
