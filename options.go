package container

type options struct {
	args      []any
	container Container
}

type Option func(*options)

func WithContainer(c Container) Option {
	return func(o *options) {
		o.container = c
	}
}

func WithArgs(args ...any) Option {
	return func(o *options) {
		o.args = args
	}
}

func buildOptions(opts ...Option) options {
	depOpts := options{
		container: get(),
	}

	for _, opt := range opts {
		opt(&depOpts)
	}

	return depOpts
}
