package httpwrap

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

type header struct {
	Integer int `http:"header=integer"`
}

type query struct {
	String string `http:"query=string"`
}

type typedContext struct {
	Integer int
	String  string
	Bool    bool
}

type typedResponse struct {
	Value int `json:"value"`
}

func TestStandardWrapper(t *testing.T) {
	var rw *httptest.ResponseRecorder
	var req *http.Request

	t.Run("simple error", func(t *testing.T) {
		wrapper := NewStandardWrapper()
		noErrorHandler := wrapper.Wrap(func() error {
			return nil
		})
		errorHandler := wrapper.Wrap(func() error {
			return NewHTTPError(http.StatusForbidden, "Forbidden.")
		})

		rw = httptest.NewRecorder()
		req = httptest.NewRequest("GET", "/endpoint", nil)
		noErrorHandler.ServeHTTP(rw, req)
		require.Equal(t, http.StatusOK, rw.Result().StatusCode)

		rw = httptest.NewRecorder()
		req = httptest.NewRequest("GET", "/endpoint", nil)
		errorHandler.ServeHTTP(rw, req)
		require.Equal(t, http.StatusForbidden, rw.Result().StatusCode)
	})

	t.Run("no middleware", func(t *testing.T) {
		wrapper := NewStandardWrapper()
		handler := wrapper.Wrap(func(p1 header, p2 query) error {
			require.Equal(t, 12, p1.Integer)
			require.Equal(t, "abc", p2.String)
			return nil
		})

		rw := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/endpoint?string=abc", nil)
		req.Header.Set("integer", "12")

		handler.ServeHTTP(rw, req)
		require.Equal(t, http.StatusOK, rw.Result().StatusCode)
	})

	t.Run("json response", func(t *testing.T) {
		wrapper := NewStandardWrapper()
		handler := wrapper.Wrap(func(p1 header, p2 query) (typedResponse, error) {
			require.Equal(t, 12, p1.Integer)
			require.Equal(t, "abc", p2.String)
			return typedResponse{Value: 42}, nil
		})

		rw := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/endpoint?string=abc", nil)
		req.Header.Set("integer", "12")

		handler.ServeHTTP(rw, req)
		statusCode, body := readResponseRecorder(t, rw)
		require.Equal(t, http.StatusOK, statusCode)
		require.Equal(t, `{"value":42}`, body)
	})

	t.Run("with middleware", func(t *testing.T) {
		wrapper := NewStandardWrapper().
			Before(func() typedContext {
				return typedContext{
					Integer: 13,
					String:  "abc",
					Bool:    true,
				}
			})

		handler := wrapper.Wrap(func(p1 header, p2 query, context typedContext) (typedResponse, error) {
			require.Equal(t, 12, p1.Integer)
			require.Equal(t, "abc", p2.String)

			require.Equal(t, 13, context.Integer)
			require.Equal(t, "abc", context.String)
			require.Equal(t, true, context.Bool)

			return typedResponse{Value: 42}, nil
		})

		rw := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/endpoint?string=abc", nil)
		req.Header.Set("integer", "12")

		handler.ServeHTTP(rw, req)
		statusCode, body := readResponseRecorder(t, rw)
		require.Equal(t, http.StatusOK, statusCode)
		require.Equal(t, `{"value":42}`, body)
	})

	t.Run("middleware shortcircuit", func(t *testing.T) {
		wrapper := NewStandardWrapper().
			Before(func() error {
				return NewHTTPError(http.StatusUnauthorized, "Unauthorized.")
			})

		handler := wrapper.Wrap(func(p1 header, p2 query, context typedContext) error {
			require.Fail(t, "Handler should never be called.")
			return nil
		})

		rw := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/endpoint?string=abc", nil)
		req.Header.Set("integer", "12")

		handler.ServeHTTP(rw, req)
		statusCode, body := readResponseRecorder(t, rw)
		require.Equal(t, http.StatusUnauthorized, statusCode)
		require.Equal(t, "Unauthorized.", body)
	})

	t.Run("middleware shortcircuit with nooperror", func(t *testing.T) {
		wrapper := NewStandardWrapper().
			Before(func(rw http.ResponseWriter) error {
				rw.WriteHeader(http.StatusCreated)
				rw.Write([]byte("HELLO WORLD"))
				return NewNoopError()
			})

		handler := wrapper.Wrap(func(p1 header, p2 query, context typedContext) error {
			require.Fail(t, "Handler should never be called.")
			return nil
		})

		rw := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/endpoint?string=abc", nil)
		req.Header.Set("integer", "12")

		handler.ServeHTTP(rw, req)
		statusCode, body := readResponseRecorder(t, rw)
		require.Equal(t, http.StatusCreated, statusCode)
		require.Equal(t, "HELLO WORLD", body)
	})
}

func readResponseRecorder(t *testing.T, rw *httptest.ResponseRecorder) (int, string) {
	result := rw.Result()
	body, err := io.ReadAll(result.Body)
	require.NoError(t, err)
	return result.StatusCode, strings.TrimSpace(string(body))
}
