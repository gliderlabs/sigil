package var_

import "reflect"

func (_ Module) Default(value, in interface{}) interface{} {
	if in == nil {
		return value
	}
	if reflect.Zero(reflect.TypeOf(in)).Interface() == in {
		return value
	}
	return in
}
