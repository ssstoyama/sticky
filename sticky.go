package sticky

type stickyContext interface {
	Value(any) any
}

// New creates a container instance.
func New(opts ...containerOption) Container {
	return newContainer(opts...)
}

// Register registers dependencies.
func Register(ctx stickyContext, rters ...register) error {
	c, err := getContainer(ctx)
	if err != nil {
		return err
	}
	for _, rter := range rters {
		if err := c.Register(rter); err != nil {
			return err
		}
	}
	return nil
}

// Resolve resolves a dependency. it can use the following options.
//
// - Tag: can resolve a dependency by T's type and tag name
func Resolve[T any](ctx stickyContext, opts ...resolveOption) (ret T, err error) {
	var c *container
	c, err = getContainer(ctx)
	if err != nil {
		return
	}
	var option resolveOptions
	for _, opt := range opts {
		opt.applyResolveOption(&option)
	}
	t := makeType[T]()
	key := dKey{t: t, tag: option.Tag}
	var v any
	v, err = c.Resolve(key)
	if err != nil {
		return
	}
	ret = v.(T)
	return
}

// Extract extracts dependencies.
func Extract(ctx stickyContext, function any) error {
	c, err := getContainer(ctx)
	if err != nil {
		return err
	}
	return c.Extract(function)
}

// Decorate allows to edit instance of generated dependencies.
func Decorate[T any](ctx stickyContext, function func(T) (T, error), opts ...resolveOption) error {
	c, err := getContainer(ctx)
	if err != nil {
		return err
	}
	var option resolveOptions
	for _, opt := range opts {
		opt.applyResolveOption(&option)
	}
	var f func(any) (any, error) = func(v any) (any, error) {
		return function(v.(T))
	}
	t := makeType[T]()
	key := dKey{t: t, tag: option.Tag}
	return c.Decorate(key, f)
}

// Validate verifies that the dependencies are registered without omission.
func Validate(ctx stickyContext) error {
	c, err := getContainer(ctx)
	if err != nil {
		return err
	}
	return c.Validate()
}
