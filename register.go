package sticky

import "reflect"

type register interface {
	Keys() ([]dKey, error)
	Scopes() ([]*dependency, error)
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

func (ini *constructorRegister) Keys() ([]dKey, error) {
	ft := reflect.TypeOf(ini.fn)
	if ft.Kind() != reflect.Func {
		return nil, &invalidConstructorError{ft}
	}
	keys := make([]dKey, ft.NumOut())
	for i := range keys {
		keys[i] = dKey{t: ft.Out(i)}
	}
	return keys, nil
}

func (ini *constructorRegister) Scopes() ([]*dependency, error) {
	fv := reflect.ValueOf(ini.fn)
	if fv.Kind() != reflect.Func {
		return nil, &invalidConstructorError{fv.Type()}
	}
	values := make([]*dependency, fv.Type().NumOut())
	for i := range values {
		values[i] = &dependency{value: fv}
	}
	return values, nil
}

func (ini *constructorRegister) Opts() []registerOption {
	return ini.opts
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

func (ini *paramRegister) Keys() ([]dKey, error) {
	keys := make([]dKey, 1)
	keys[0] = dKey{
		t:   reflect.TypeOf(ini.value),
		tag: ini.tag,
	}
	return keys, nil
}

func (ini *paramRegister) Scopes() ([]*dependency, error) {
	values := make([]*dependency, 1)
	values[0] = &dependency{
		value:   reflect.ValueOf(ini.value),
		isParam: true,
	}
	return values, nil
}

func (ini *paramRegister) Opts() []registerOption {
	return make([]registerOption, 0)
}
