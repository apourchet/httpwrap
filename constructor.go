package httpwrap

import "net/http"

// Constructor is the function signature for unmarshalling an http request into
// an object.
type Constructor func(http.ResponseWriter, *http.Request, any) error

// emptyConstructor is the default constructor for new wrappers.
// It is a no-op, and will not parse any http request information to construct endpoint
// parameter objects.
func emptyConstructor(http.ResponseWriter, *http.Request, any) error { return nil }
