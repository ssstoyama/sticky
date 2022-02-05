package sticky

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
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

type cycleDependencyError struct {
	deps []reflect.Type
}

func (e *cycleDependencyError) Error() string {
	deps := make([]string, 0, len(e.deps))
	for i := len(e.deps) - 1; i >= 0; i-- {
		deps = append(deps, fmt.Sprintf("%s%s", strings.Repeat(" ", len(e.deps)-i-1), pathString(e.deps[i])))
	}
	return fmt.Sprintf("cycle dependency error.\n%s", strings.Join(deps, "\n"))
}
