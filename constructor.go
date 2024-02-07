package httpwrap

import "net/http"

// RequestReader is the function signature for unmarshalling a *http.Request into
// an object.
type RequestReader func(http.ResponseWriter, *http.Request, any) error

// ResponseWriter is the function signature for marshalling a structured response
// into a standard http.ResponseWriter
type ResponseWriter func(w http.ResponseWriter, r *http.Request, res any, err error)

// emptyRequestReader is the default constructor for new wrappers.
// It is a no-op, and will not parse any http request information to construct endpoint
// parameter objects.
func emptyRequestReader(http.ResponseWriter, *http.Request, any) error { return nil }
