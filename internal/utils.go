package internal

import "reflect"

// DerefStruct dereferences the type of the object if it is
// a pointer or an interface.
// Returns whether the final type is a struct.
func DerefStruct(obj interface{}) (reflect.Type, bool) {
	st := reflect.TypeOf(obj)
	for st.Kind() == reflect.Ptr || st.Kind() == reflect.Interface {
		st = st.Elem()
	}

	return st, st.Kind() == reflect.Struct
}
