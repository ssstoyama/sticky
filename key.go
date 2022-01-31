package sticky

import "reflect"

type stickey string

const defaultKey stickey = "fingers"

type dKey struct {
	t   reflect.Type
	tag string
}

func (k dKey) Type() reflect.Type {
	return k.t
}

func (k dKey) Tag() string {
	return k.tag
}

func (k dKey) IsInterfaceType() bool {
	return k.t.Kind() == reflect.Interface
}

func (k dKey) IsErrorType() bool {
	errT := reflect.TypeOf((*error)(nil)).Elem()
	return k.t.Implements(errT)
}

func (s *dKey) applyOption(opt *registerOptions) error {
	if s.tag == "" {
		s.tag = opt.Tag
	}
	if opt.Implements != nil {
		s.t = *opt.Implements
	}
	return nil
}
