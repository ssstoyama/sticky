package sticky

import (
	"errors"
	"reflect"
)

type dependency struct {
	value      reflect.Value
	instance   any
	implements *reflect.Type
	isParam    bool
	cache      *bool
}

func (s *dependency) getValue() (any, bool) {
	return s.instance, s.instance != nil
}

func (s *dependency) returnTypes() []reflect.Type {
	out := make([]reflect.Type, s.value.Type().NumOut())
	for i := range out {
		out[i] = s.value.Type().Out(i)
	}
	return out
}

func (s *dependency) applyOption(opt *registerOptions) error {
	s.cache = opt.Cache

	if opt.Implements == nil {
		return nil
	}
	it := *opt.Implements
	rt := s.returnTypes()[0]
	if !rt.Implements(it) {
		return errors.New("not implements")
	}
	s.implements = opt.Implements
	return nil
}
