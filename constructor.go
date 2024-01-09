package httpwrap

import "net/http"

// Constructor is the function signature for unmarshalling an http request into
// an object.
type Constructor func(http.ResponseWriter, *http.Request, any) error

// EmptyConstructor is the default constructor for new wrappers.
// It is a no-op.
func EmptyConstructor(http.ResponseWriter, *http.Request, any) error { return nil }
