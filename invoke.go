package sticky

import "reflect"

type invoker func(reflect.Value, []reflect.Value) []reflect.Value

func defaultInvoker(fn reflect.Value, args []reflect.Value) []reflect.Value {
	return fn.Call(args)
}

func dryInvoker(fn reflect.Value, _ []reflect.Value) []reflect.Value {
	ft := fn.Type()
	results := make([]reflect.Value, ft.NumOut())
	for i := range results {
		results[i] = reflect.Zero(fn.Type().Out(i))
	}
	return results
}
