package httpwrap

import (
	"encoding/json"
	"log"
	"net/http"
)

// The StandardConstructor decodes the request using the following:
// - cookies
// - query params
// - path params
// - headers
// - JSON decoding of the body
func StandardConstructor() Constructor {
	decoder := NewDecoder()
	return func(rw http.ResponseWriter, req *http.Request, obj any) error {
		return decoder.Decode(req, obj)
	}
}

// StandardResponseWriter will try to cast the error and response objects to the
// HTTPResponse interface and use them to send the response to the client.
// By default, it will send a 200 OK and encode the response object as JSON.
// If the HTTPResponse has a `0` StatusCode, WriteHeader will not be called.
// If the error is not an HTTPResponse, a 500 status code will be returned with
// the body being exactly the error's string.
func StandardResponseWriter() func(w http.ResponseWriter, res any, err error) {
	return func(w http.ResponseWriter, res any, err error) {
		if err != nil {
			if cast, ok := err.(HTTPResponse); ok {
				code := cast.StatusCode()
				if code != 0 {
					w.WriteHeader(cast.StatusCode())
				}
				if sendError := cast.WriteBody(w); sendError != nil {
					log.Println("error writing response:", sendError)
				}
			} else {
				w.WriteHeader(http.StatusInternalServerError)
				if _, sendError := w.Write([]byte(err.Error())); sendError != nil {
					log.Println("error writing response:", sendError)
				}
			}
			return
		}

		if res == nil {
			return
		}

		if cast, ok := res.(HTTPResponse); ok {
			code := cast.StatusCode()
			if code != 0 {
				w.WriteHeader(cast.StatusCode())
			}
			if sendError := cast.WriteBody(w); sendError != nil {
				log.Println("error writing response:", sendError)
			}
			return
		}

		w.WriteHeader(http.StatusOK)
		encoder := json.NewEncoder(w)
		encoder.SetIndent("", "  ")
		if sendError := encoder.Encode(res); sendError != nil {
			log.Println("Error writing response:", sendError)
		}
	}
}

// NewStandardWrapper returns a new wrapper using the StandardConstructor and the
// StandardResponseWriter.
func NewStandardWrapper() Wrapper {
	constructor := StandardConstructor()
	responseWriter := StandardResponseWriter()
	return New().
		WithConstruct(constructor).
		Finally(responseWriter)
}
