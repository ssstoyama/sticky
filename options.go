package sticky

import "reflect"

// containerOption is interface to apply option.
type containerOption interface {
	applyContainerOption(*containerOptions)
}

// containerOptions is for the container.
type containerOptions struct {
	Cache bool
}

// registerOption is interface to apply option.
type registerOption interface {
	applyRegisterOption(*registerOptions)
}

// registerOptions is for the Register method.
type registerOptions struct {
	Tag        string
	Implements *reflect.Type
	Cache      *bool
}

// resolveOption is interface to apply option.
type resolveOption interface {
	applyResolveOption(*resolveOptions)
}

// resolveOptions is for the Resolve method.
type resolveOptions struct {
	Tag string
}

// Tag option allows to tag dependencies.
//
// e.g.
// - Register(c, Constructor(/* some constructor */, Tag("MyTag")))
// - Resolve[T](c, Tag("MyTag"))
func Tag(tag string) *tagOption {
	return &tagOption{tag: tag}
}

type tagOption struct {
	tag string
}

func (o *tagOption) applyRegisterOption(opt *registerOptions) {
	opt.Tag = o.tag
}

func (o *tagOption) applyResolveOption(opt *resolveOptions) {
	opt.Tag = o.tag
}

// Implements option allows to register dependencies as interfaces.
//
// e.g.
// - Register(c, Constructor(/* some constructor */, Implements[InterfaceType]()))
func Implements[T any]() *implementsOption {
	return &implementsOption{t: makeType[T]()}
}

type implementsOption struct{ t reflect.Type }

func (o *implementsOption) applyRegisterOption(opt *registerOptions) {
	opt.Implements = &o.t
}

// Cache option can be used to reuse the generated dependencies.
//
// e.g.
// - New(Cache(true)) // default true
// - Register(c, Constructor(/* some constructor */, Cache(true)))
func Cache(enable bool) *cacheOption {
	return &cacheOption{enable}
}

type cacheOption struct{ enable bool }

func (o *cacheOption) applyContainerOption(opt *containerOptions) {
	opt.Cache = o.enable
}

func (o *cacheOption) applyRegisterOption(opt *registerOptions) {
	opt.Cache = &o.enable
}
