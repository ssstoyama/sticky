package sticky

import (
	"errors"
	"fmt"
	"reflect"
)

// assertConstructor is determine if v is a constructor.
// constructor must be Function and returns value.
func assertConstructor(v reflect.Value) error {
	if v.Kind() != reflect.Func {
		return &invalidFunctionError{}
	}
	t := v.Type()
	if t.NumOut() < 1 {
		return &invalidConstructorError{t}
	}
	return nil
}

// assertNotCycle makes sure that the dependencies are not cycle.
func assertNotCycle(c *container, deps []*dependency) error {
	for _, dep := range deps {
		t := dep.value.Type()
		if t.Kind() != reflect.Func {
			continue
		}
		var cerr cycleDependencyError
		for i := 0; i < t.NumOut(); i++ {
			cerr.deps = append(cerr.deps, t.Out(i))
		}
		if err := _assertNotCycle(c, dep, cerr); err != nil {
			return err
		}
	}
	return nil
}

func _assertNotCycle(c *container, dep *dependency, cerr cycleDependencyError) error {
	t := dep.value.Type()
	for i := 0; i < t.NumIn(); i++ {
		it := t.In(i)
		for _, dt := range cerr.deps {
			if it == dt {
				cerr.deps = append(cerr.deps, it)
				return &cerr
			}
		}
		cerr.deps = append(cerr.deps, it)
		dep, ok := c.dependencies[dKey{t: it}]
		if !ok {
			continue
		}
		if err := _assertNotCycle(c, dep, cerr); err != nil {
			return err
		}
	}
	return nil
}

// get constructor from context.
func getContainer(ctx stickyContext) (*container, error) {
	if c, ok := ctx.(*container); ok {
		return c, nil
	}
	c := ctx.Value(defaultKey)
	if c == nil {
		return nil, errors.New("not found container in context")
	}
	return c.(*container), nil
}

// indirectType returns the type that t points to
func indirectType(t reflect.Type) reflect.Type {
	if t.Kind() == reflect.Pointer {
		return t.Elem()
	}
	return t
}

// make reflect.Type from T.
func makeType[T any]() reflect.Type {
	return reflect.TypeOf((*T)(nil)).Elem()
}

func pathString(t reflect.Type) string {
	_t := indirectType(t)
	path := fmt.Sprintf("%s.%s", _t.PkgPath(), _t.Name())
	if t.Kind() == reflect.Ptr {
		return "*" + path
	}
	return path
}
