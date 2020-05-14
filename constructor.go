package httpwrap

import "net/http"

// Constructor is the function signature for unmarshalling an http request into
// an object.
type Constructor func(http.ResponseWriter, *http.Request, interface{}) error

// EmptyConstructor is the default constructor for new wrappers.
// It is a no-op.
func EmptyConstructor(http.ResponseWriter, *http.Request, interface{}) error { return nil }

// The StandardConstructor decodes the request using the following:
// - cookies
// - query params
// - path params
// - headers
// - JSON decoding of the body
func StandardConstructor() Constructor {
	decoder := NewDecoder()
	return func(rw http.ResponseWriter, req *http.Request, obj interface{}) error {
		return decoder.Decode(req, obj)
	}
}
