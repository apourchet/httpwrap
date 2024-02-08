package httpwrap

import (
	"encoding/json"
	"io"
)

// HTTPResponse is used by the StandardResponseWriter to construct the
// response body according to the StatusCode() and WriteBody() functions.
// If the StatusCode() function returns `0`, the StandardResponseWriter will
// assume that WriteHeader has already been called on the http.ResponseWriter
// object.
type HTTPResponse interface {
	StatusCode() int
	WriteBody(io.Writer) error
}

// The jsonResponse type implements HTTPResponse. When returned, it will
// write the status code in the http response's header and JSON encode the
// body.
type jsonResponse struct {
	code int
	body any
}

func NewJSONResponse(code int, body any) HTTPResponse {
	return jsonResponse{
		code: code,
		body: body,
	}
}

func (res jsonResponse) StatusCode() int { return res.code }

func (res jsonResponse) WriteBody(writer io.Writer) error {
	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", "  ")
	return encoder.Encode(res.body)
}

type noopResponse struct{}

// NewNoopResponse returns an HTTPResponse that will completely bypass the
// deserialization logic. This can be used when the endpoint or middleware
// operates directly on the native http.ResponseWriter.
func NewNoopResponse() HTTPResponse {
	return &noopResponse{}
}

func (res noopResponse) StatusCode() int { return 0 }

func (res noopResponse) WriteBody(writer io.Writer) error { return nil }
