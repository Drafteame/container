package utils

import (
	"reflect"
)

// GetFirstReturnType Receive a callback as an input and try to return the first return type as `reflect.Type` element.
// Return nil of the construct provided is not a function, is nil or do not return anything.
func GetFirstReturnType(construct any) reflect.Type {
	ctype := reflect.TypeOf(construct)

	if ctype == nil || ctype.Kind() != reflect.Func || ctype.NumOut() < 1 {
		return nil
	}

	return ctype.Out(0)
}
