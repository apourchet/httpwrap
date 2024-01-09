package httpwrap

import (
	"fmt"
	"io"
)

// HTTPError implements both the HTTPResponse interface and the standard error
// interface.
type HTTPError struct {
	code int
	body string
}

func NewHTTPError(code int, format string, args ...any) HTTPError {
	return HTTPError{
		code: code,
		body: fmt.Sprintf(format, args...),
	}
}

func (err HTTPError) Error() string {
	return fmt.Sprintf("http error: %d: %s", err.code, err.body)
}

func (err HTTPError) StatusCode() int { return err.code }

func (err HTTPError) WriteBody(writer io.Writer) error {
	_, writeError := io.WriteString(writer, err.body)
	return writeError
}
