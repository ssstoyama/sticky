package sticky

import (
	"bytes"
	"fmt"
	"reflect"
)

type alreadyRegisteredError struct {
	key dKey
}

func (e *alreadyRegisteredError) Error() string {
	tag := e.key.tag
	if tag == "" {
		tag = `''`
	}
	return fmt.Sprintf("already registered: type=%s, tag=%s", pathString(e.key.t), tag)
}

type notFoundRegisterError struct {
	key dKey
}

func (e *notFoundRegisterError) Error() string {
	tag := e.key.tag
	if tag == "" {
		tag = `''`
	}
	return fmt.Sprintf("not found register: type=%s, tag=%s", pathString(e.key.t), tag)
}

type invalidFunctionError struct {
}

func (e *invalidFunctionError) Error() string {
	return "invalid value. must be function"
}

type invalidConstructorError struct {
	t reflect.Type
}

func (e *invalidConstructorError) Error() string {
	return fmt.Sprintf("invalid constructor. must be factory function. got=%s", e.t.Kind())
}

type validationError struct {
	errs []error
}

func (e *validationError) Error() string {
	var buf bytes.Buffer
	buf.WriteString("validation error:")
	for _, err := range e.errs {
		buf.WriteString(fmt.Sprintf("\n\t%s", err.Error()))
	}
	return buf.String()
}

func (e *validationError) IsError() bool {
	return len(e.errs) > 0
}
