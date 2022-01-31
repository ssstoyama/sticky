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
