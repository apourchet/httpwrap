package httpwrap

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWrapper(t *testing.T) {
	t.Run("just main", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test", nil)
		rw := httptest.NewRecorder()

		handler := New().
			WithConstruct(nopConstructor).
			Wrap(func(rw http.ResponseWriter, req *http.Request) error {
				require.NotNil(t, rw)
				require.Equal(t, "GET", req.Method)
				return nil
			})
		handler.ServeHTTP(rw, req)
		require.Equal(t, http.StatusOK, rw.Result().StatusCode)
	})

	t.Run("main with after", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test", nil)
		rw := httptest.NewRecorder()

		type resp struct{ s string }
		handler := New().
			WithConstruct(nopConstructor).
			Finally(func(res interface{}, err error) {
				s := fmt.Sprintf("%v", res)
				require.True(t, strings.Contains(s, "response"))
				require.Error(t, err)
			}).
			Wrap(func() (resp, error) {
				return resp{"response"}, fmt.Errorf("error")
			})
		handler.ServeHTTP(rw, req)
	})

	t.Run("main with before", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test", nil)
		rw := httptest.NewRecorder()

		type meta struct{ path string }
		handler := New().
			WithConstruct(nopConstructor).
			Before(func(req *http.Request) meta {
				return meta{req.URL.Path}
			}).
			Wrap(func(m meta) {
				require.Equal(t, "/test", m.path)
			})
		handler.ServeHTTP(rw, req)
	})

	t.Run("before failing", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test", nil)
		rw := httptest.NewRecorder()

		type meta struct{ path string }
		handler := New().
			WithConstruct(nopConstructor).
			Before(func(req *http.Request) (meta, error) {
				return meta{req.URL.Path}, fmt.Errorf("failed before")
			}).
			Finally(func(rw http.ResponseWriter, m meta, res interface{}, err error) {
				require.Equal(t, "/test", m.path)
				require.Nil(t, res)
				require.Error(t, err)
			}).
			Wrap(func(m meta) error {
				require.FailNow(t, "should not call main handler")
				return nil
			})
		handler.ServeHTTP(rw, req)
	})

	t.Run("main failing", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test", nil)
		rw := httptest.NewRecorder()

		type resp struct{}
		type meta struct{ path string }
		handler := New().
			WithConstruct(nopConstructor).
			Before(func(req *http.Request) (meta, error) {
				return meta{req.URL.Path}, nil
			}).
			Finally(func(rw http.ResponseWriter, m meta, res interface{}, err error) {
				require.Equal(t, "/test", m.path)
				require.Nil(t, res)
				require.Error(t, err)
			}).
			Wrap(func(m meta) (*resp, error) {
				return nil, fmt.Errorf("main failed")
			})
		handler.ServeHTTP(rw, req)
	})
}
