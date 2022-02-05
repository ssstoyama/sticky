package sticky

import (
	"context"
	"errors"
	"reflect"
)

// Container is DI container
type Container interface {
	stickyContext
	WithContext(ctx context.Context) context.Context
}

func newContainer(opts ...containerOption) *container {
	c := &container{
		dependencies: make(map[dKey]*dependency),
		cache:        true,
		invoker:      defaultInvoker,
	}
	option := containerOptions{
		Cache: true,
	}
	for _, opt := range opts {
		opt.applyContainerOption(&option)
	}
	c.cache = option.Cache
	return c
}

type container struct {
	dependencies map[dKey]*dependency
	cache        bool
	invoker      invoker
}

// WithContext saves the container in the context and returns it.
// context in which the container is saved can be passed as an argument to sticky.
func (c *container) WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, defaultKey, c)
}

// Value is a method to be implemented to satisfy the stickyContext interface.
func (c *container) Value(key any) any {
	return errors.New("not implemented")
}

// Register registers a dependency.
func (c *container) Register(ini register) error {
	keys, err := ini.Keys()
	if err != nil {
		return err
	}
	dependencies, err := ini.Scopes()
	if err != nil {
		return err
	}
	if err := assertNotCycle(c, dependencies); err != nil {
		return err
	}
	opts := ini.Opts()

	for i := range keys {
		key := keys[i]
		dep := dependencies[i]

		if key.IsErrorType() {
			continue
		}

		options := registerOptions{}
		for _, opt := range opts {
			opt.applyRegisterOption(&options)
		}
		if err := c.applyRegisterOption(&key, dep, &options); err != nil {
			return err
		}

		if !dep.isParam {
			if err := assertConstructor(dep.value); err != nil {
				return err
			}
		}

		if _, ok := c.dependencies[key]; ok {
			return &alreadyRegisteredError{key}
		}
		c.dependencies[key] = dep
	}
	return nil
}

// Resolve resolves a dependency.
func (c *container) Resolve(key dKey) (any, error) {
	dep, err := c.findDep(key)
	if err != nil {
		return nil, err
	}
	if dep.isParam {
		return dep.value.Interface(), nil
	}
	if v, ok := dep.getValue(); ok {
		return v, nil
	}
	values, err := c.call(dep.value)
	if err != nil {
		return nil, err
	}
	result, err := c.pick(key.t, values)
	if err != nil {
		return nil, err
	}
	if err := c.commit(key, values); err != nil {
		return nil, err
	}
	return result, nil
}

// Extract extracts dependency.
func (c *container) Extract(function any) error {
	fnV := reflect.ValueOf(function)
	if fnV.Kind() != reflect.Func {
		return &invalidFunctionError{}
	}

	fnT := fnV.Type()
	args := make([]reflect.Value, fnT.NumIn())
	for i := range args {
		inT := fnT.In(i)
		arg, err := c.Resolve(dKey{t: inT})
		if err != nil {
			return err
		}
		args[i] = reflect.ValueOf(arg)
	}
	fnV.Call(args)
	return nil
}

// Decorate allows to edit instance of generated dependencies.
func (c *container) Decorate(key dKey, function func(any) (any, error)) error {
	dep, err := c.findDep(key)
	if err != nil {
		return err
	}
	v, err := c.Resolve(key)
	if err != nil {
		return err
	}
	decorated, err := function(v)
	if err != nil {
		return err
	}
	dep.instance = decorated
	var cached = true
	dep.cache = &cached
	return nil
}

// Validate verifies that the dependencies are registered without omission.
func (c *container) Validate() error {
	_c := *c
	_c.cache = false
	_c.invoker = dryInvoker
	var vErr validationError
	for _, dep := range c.dependencies {
		if dep.isParam {
			continue
		}
		fn := dep.value
		if _, err := _c.call(fn); err != nil {
			vErr.errs = append(vErr.errs, err)
		}
	}
	if vErr.IsError() {
		return &vErr
	}
	return nil
}

func (c *container) applyRegisterOption(key *dKey, dep *dependency, options *registerOptions) error {
	if err := key.applyOption(options); err != nil {
		return err
	}
	if err := dep.applyOption(options); err != nil {
		return err
	}
	return nil
}

func (c *container) findDep(key dKey) (*dependency, error) {
	dep, ok := c.dependencies[key]
	if !ok {
		return nil, &notFoundRegisterError{key}
	}
	return dep, nil
}

// pick returns a value that matches t's type in values.
// if value implements error interface, it is converted to error type and returns.
func (c *container) pick(t reflect.Type, values []any) (ret any, err error) {
	for _, v := range values {
		if v == nil {
			continue
		}
		if reflect.TypeOf(v).AssignableTo(t) {
			ret = v
		}
		if er, ok := v.(error); ok {
			err = er
		}
	}
	return
}

// call returns result of executing the constructor function.
func (c *container) call(fn reflect.Value) ([]any, error) {
	fnT := fn.Type()
	args := make([]reflect.Value, fnT.NumIn())
	for i := range args {
		inT := fnT.In(i)
		arg, err := c.Resolve(dKey{t: inT})
		if err != nil {
			return nil, err
		}
		args[i] = reflect.ValueOf(arg)
	}
	results := make([]any, fn.Type().NumOut())
	for i, ret := range c.invoker(fn, args) {
		results[i] = ret.Interface()
	}
	return results, nil
}

// commit stores generated dependencies if cache option is enabled.
func (c *container) commit(sKey dKey, values []any) error {
	for _, value := range values {
		if value == nil {
			continue
		}
		key := dKey{t: reflect.TypeOf(value), tag: sKey.tag}
		if sKey.IsInterfaceType() {
			key.t = sKey.Type()
		}
		dep, err := c.findDep(key)
		if err != nil {
			return err
		}
		if dep.cache == nil {
			if c.cache {
				dep.instance = value
			}
			continue
		}
		if *dep.cache {
			dep.instance = value
		}
	}
	return nil
}
