package httpwrap

import (
	"fmt"
	"io"
)

type HTTPError interface {
	error
	HTTPResponse
}

// httpError implements both the HTTPResponse interface and the standard error
// interface.
type httpError struct {
	code int
	body string
}

func NewHTTPError(code int, format string, args ...any) HTTPError {
	return httpError{
		code: code,
		body: fmt.Sprintf(format, args...),
	}
}

func (err httpError) Error() string {
	return fmt.Sprintf("http error: %d: %s", err.code, err.body)
}

func (err httpError) StatusCode() int { return err.code }

func (err httpError) WriteBody(writer io.Writer) error {
	_, writeError := io.WriteString(writer, err.body)
	return writeError
}

// NewNoopError returns an HTTPError that will completely bypass the
// deserialization logic. This can be used when the endpoint or middleware
// operates directly on the native http.ResponseWriter.
func NewNoopError() HTTPError { return NewHTTPError(0, "") }
