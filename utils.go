package httpwrap

import (
	"net/http"
	"reflect"
)

var _httpResponseWriterType = reflect.TypeOf(http.ResponseWriter(nil))
var _httpRequestType = reflect.TypeOf(&http.Request{})
var _errorType = reflect.TypeOf((*error)(nil))

func isEmptyInterface(t reflect.Type) bool {
	return t.String() == "interface {}"
}

func isError(t reflect.Type) bool {
	return t.String() == "error"
}

func validateBefore(in, out []reflect.Type) error {
	// TODO: number of intypes that are emptyInterface == 0
	// TODO: number of intypes that are error <= 1
	return nil
}

func validateMain(in, out []reflect.Type) error {
	// TODO: Assert that input types is never interface.
	// TODO: Assert that first output type isnt error if len(outs) >= 2.
	return nil
}

func validateAfter(in, out []reflect.Type) error {
	// TODO: number of intypes that are emptyInterface <= 1
	// TODO: number of intypes that are error <= 1
	return nil
}
