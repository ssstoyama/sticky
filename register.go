package sticky

import "reflect"

type register interface {
	Keys() ([]dKey, error)
	Deps() ([]*dependency, error)
	Opts() []registerOption
}

func Constructor(fn any, opts ...registerOption) *constructorRegister {
	return &constructorRegister{
		fn:   fn,
		opts: opts,
	}
}

type constructorRegister struct {
	fn   any
	opts []registerOption
}

func (cr *constructorRegister) Keys() ([]dKey, error) {
	ft := reflect.TypeOf(cr.fn)
	if ft.Kind() != reflect.Func {
		return nil, &invalidConstructorError{ft}
	}
	keys := make([]dKey, ft.NumOut())
	for i := range keys {
		keys[i] = dKey{t: ft.Out(i)}
	}
	return keys, nil
}

func (cr *constructorRegister) Deps() ([]*dependency, error) {
	fv := reflect.ValueOf(cr.fn)
	if fv.Kind() != reflect.Func {
		return nil, &invalidConstructorError{fv.Type()}
	}
	values := make([]*dependency, fv.Type().NumOut())
	for i := range values {
		values[i] = &dependency{value: fv}
	}
	return values, nil
}

func (cr *constructorRegister) Opts() []registerOption {
	return cr.opts
}

func Param(value any, tag string) *paramRegister {
	return &paramRegister{
		value: value,
		tag:   tag,
	}
}

type paramRegister struct {
	tag   string
	value any
}

func (pr *paramRegister) Keys() ([]dKey, error) {
	keys := make([]dKey, 1)
	keys[0] = dKey{
		t:   reflect.TypeOf(pr.value),
		tag: pr.tag,
	}
	return keys, nil
}

func (pr *paramRegister) Deps() ([]*dependency, error) {
	values := make([]*dependency, 1)
	values[0] = &dependency{
		value:   reflect.ValueOf(pr.value),
		isParam: true,
	}
	return values, nil
}

func (pr *paramRegister) Opts() []registerOption {
	return make([]registerOption, 0)
}
