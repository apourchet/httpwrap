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

// The JSONResponse type implements HTTPResponse. When returned, it will
// write the status code in the http response's header and JSON encode the
// body.
type JSONResponse struct {
	code int
	body any
}

func NewJSONResponse(code int, body any) HTTPResponse {
	return JSONResponse{
		code: code,
		body: body,
	}
}

func (res JSONResponse) StatusCode() int { return res.code }

func (res JSONResponse) WriteBody(writer io.Writer) error {
	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", "  ")
	return encoder.Encode(res.body)
}
