package httpwrap

import (
	"fmt"
	"reflect"
)

var _errorType = reflect.TypeOf((*error)(nil)).Elem()

func isEmptyInterface(t reflect.Type) bool {
	return t.String() == "interface {}"
}

func isError(t reflect.Type) bool {
	return t.Implements(_errorType)
}

func typesOf(fn any) ([]reflect.Type, []reflect.Type) {
	val := reflect.ValueOf(fn)
	fnType := val.Type()
	inTypes, outTypes := []reflect.Type{}, []reflect.Type{}
	for i := 0; i < fnType.NumIn(); i++ {
		inTypes = append(inTypes, fnType.In(i))
	}
	for i := 0; i < fnType.NumOut(); i++ {
		outTypes = append(outTypes, fnType.Out(i))
	}
	return inTypes, outTypes
}

func validateBefore(in, _ []reflect.Type) error {
	for i, t := range in {
		if isEmptyInterface(t) {
			return fmt.Errorf("before input #%d must not be empty interface", i)
		}
	}
	if err := areTypesUnique(in); err != nil {
		return fmt.Errorf("before input types must be unique: %v", err)
	}
	return nil
}

func validateMain(in, _ []reflect.Type) error {
	for i, t := range in {
		if isEmptyInterface(t) {
			return fmt.Errorf("main input #%d must not be empty interface", i)
		}
	}
	if err := areTypesUnique(in); err != nil {
		return fmt.Errorf("main input types must be unique: %v", err)
	}
	// TODO: Assert check that main returns a non-error output?
	return nil
}

func validateAfter(in, _ []reflect.Type) error {
	if err := areTypesUnique(in); err != nil {
		return fmt.Errorf("after input types must be unique: %v", err)
	}
	return nil
}

func areTypesUnique(ts []reflect.Type) error {
	m := map[reflect.Type]int{}
	for i, t := range ts {
		if j, found := m[t]; found {
			return fmt.Errorf("types %d and %d are equal", j, i)
		}
		m[t] = i
	}
	return nil
}
