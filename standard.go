package httpwrap

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type HTTPResponse interface {
	StatusCode() int
	WriteBody(io.Writer) error
}

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
func StandardResponseWriter() func(w http.ResponseWriter, res any, err error) {
	return func(w http.ResponseWriter, res any, err error) {
		if err != nil {
			if cast, ok := err.(HTTPResponse); ok {
				w.WriteHeader(cast.StatusCode())
				if sendError := cast.WriteBody(w); sendError != nil {
					log.Println("error writing response:", sendError)
				}
			} else {
				w.WriteHeader(http.StatusInternalServerError)
				if _, sendError := w.Write([]byte(err.Error() + "\n")); sendError != nil {
					log.Println("error writing response:", sendError)
				}
			}
			return
		}

		if cast, ok := res.(HTTPResponse); ok {
			w.WriteHeader(cast.StatusCode())
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
