package httpwrap

import (
	"encoding/json"
	"log"
	"net/http"
)

// The StandardRequestReader decodes the request using the following:
//
// - Cookies
//
// - Query Params
//
// - Request Headers
//
// - Request Path Segment (e.g: /api/pets/{id})
//
// - JSON Decoding of the http request body
func StandardRequestReader() RequestReader {
	decoder := NewDecoder()
	return func(_ http.ResponseWriter, req *http.Request, obj any) error {
		return decoder.Decode(req, obj)
	}
}

// StandardResponseWriter will try to cast the error and response objects to the
// HTTPResponse interface and use them to send the response to the client.
// By default, it will send a 200 OK and encode the response object as JSON.
// If the HTTPResponse has a `0` StatusCode, WriteHeader will not be called.
// If the error is not an HTTPResponse, a 500 status code will be returned with
// the body being exactly the error's string.
func StandardResponseWriter() ResponseWriter {
	return func(w http.ResponseWriter, _ *http.Request, res any, err error) {
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
		w.Header().Set("Content-Type", "application/json")
		encoder := json.NewEncoder(w)
		encoder.SetEscapeHTML(false)
		if sendError := encoder.Encode(res); sendError != nil {
			log.Println("Error writing response:", sendError)
		}
	}
}

// NewStandardWrapper returns a new wrapper using the StandardRequestReader and the
// StandardResponseWriter.
func NewStandardWrapper() Wrapper {
	constructor := StandardRequestReader()
	responseWriter := StandardResponseWriter()
	return New().
		WithRequestReader(constructor).
		Finally(responseWriter)
}
