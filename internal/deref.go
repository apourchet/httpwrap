package internal

import "reflect"

// DerefType dereferences the type of the object if it is
// a pointer or an interface.
// Returns whether the final type is a struct.
func DerefType(obj interface{}) (reflect.Type, bool) {
	st := reflect.TypeOf(obj)
	for st.Kind() == reflect.Ptr || st.Kind() == reflect.Interface {
		st = st.Elem()
	}

	return st, st.Kind() == reflect.Struct
}

// DerefValue dereferences the value of the object until
// it is no longer a pointer or an interface. Also returns
// false if the underlying value is Nil.
func DerefValue(obj interface{}) (reflect.Value, bool) {
	v := reflect.ValueOf(obj)
	for v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		if v.IsNil() { // If the chain ends in a nil, skip this
			return v, false
		}
		v = v.Elem()
	}
	return v, true
}
