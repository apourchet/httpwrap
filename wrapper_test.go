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
				rw.WriteHeader(http.StatusCreated)
			}).
			Wrap(func() (resp, error) {
				return resp{"response"}, fmt.Errorf("error")
			})
		handler.ServeHTTP(rw, req)
		require.Equal(t, http.StatusCreated, rw.Result().StatusCode)
	})

	t.Run("main with before", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test", nil)
		rw := httptest.NewRecorder()

		type meta struct{ path string }
		handler := New().
			WithConstruct(nopConstructor).
			Before(func(req *http.Request) meta {
				require.NotNil(t, req)
				require.NotNil(t, req.URL)
				return meta{req.URL.Path}
			}).
			Wrap(func(m meta) {
				require.Equal(t, "/test", m.path)
				rw.WriteHeader(http.StatusCreated)
			})
		handler.ServeHTTP(rw, req)
		require.Equal(t, http.StatusCreated, rw.Result().StatusCode)
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
				require.NotNil(t, rw)
				require.Equal(t, "/test", m.path)
				require.Nil(t, res)
				require.Error(t, err)
				rw.WriteHeader(http.StatusCreated)
			}).
			Wrap(func(m meta) error {
				require.FailNow(t, "should not call main handler")
				return nil
			})
		handler.ServeHTTP(rw, req)
		require.Equal(t, http.StatusCreated, rw.Result().StatusCode)
	})

	t.Run("main failing", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test", nil)
		rw := httptest.NewRecorder()

		type resp struct{}
		type meta struct{ path string }
		handler := New().
			WithConstruct(nopConstructor).
			Before(func(req *http.Request) (meta, error) {
				require.NotNil(t, req)
				require.NotNil(t, req.URL)
				return meta{req.URL.Path}, nil
			}).
			Finally(func(rw http.ResponseWriter, m meta, res interface{}, err error) {
				require.NotNil(t, rw)
				require.Equal(t, "/test", m.path)
				require.Nil(t, res)
				require.Error(t, err)
				rw.WriteHeader(http.StatusCreated)
			}).
			Wrap(func(m meta) (*resp, error) {
				return nil, fmt.Errorf("main failed")
			})
		handler.ServeHTTP(rw, req)
		require.Equal(t, http.StatusCreated, rw.Result().StatusCode)
	})

	t.Run("with failed constructor", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test", nil)
		rw := httptest.NewRecorder()

		type meta struct{}
		handler := New().
			WithConstruct(failedConstructor).
			Before(func(m meta) {
				require.FailNow(t, "should not get to before")
			}).
			Finally(func(rw http.ResponseWriter, err error) {
				require.Error(t, err)
				require.NotNil(t, rw)
				rw.WriteHeader(http.StatusCreated)
			}).
			Wrap(func() {
				require.FailNow(t, "should not get to before")
			})
		handler.ServeHTTP(rw, req)
		require.Equal(t, http.StatusCreated, rw.Result().StatusCode)
	})

	t.Run("with constructor", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test", strings.NewReader(`{"metafield": "metafield", "field2": "field2"}`))
		rw := httptest.NewRecorder()

		type meta struct{ Metafield string }
		type extra struct{ Field1 string }
		type mainArg struct{ Field2 string }
		handler := New().
			WithConstruct(jsonBodyConstructor).
			Before(func(m meta) extra {
				require.Equal(t, "metafield", m.Metafield)
				return extra{"field1"}
			}).
			Finally(func(rw http.ResponseWriter, err error) {
				require.NoError(t, err)
				require.NotNil(t, rw)
				rw.WriteHeader(http.StatusCreated)
			}).
			Wrap(func(e extra, m mainArg) error {
				require.Equal(t, "field1", e.Field1)
				require.Equal(t, "field2", m.Field2)
				return nil
			})
		handler.ServeHTTP(rw, req)
		require.Equal(t, http.StatusCreated, rw.Result().StatusCode)
	})

	t.Run("before returns special error", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test", nil)
		rw := httptest.NewRecorder()

		handler := New().
			WithConstruct(nopConstructor).
			Before(func() *myerr {
				return &myerr{}
			}).
			Finally(func(rw http.ResponseWriter, res interface{}, err error) {
				require.NotNil(t, rw)
				require.Nil(t, res)
				require.Error(t, err)
				rw.WriteHeader(http.StatusCreated)
			}).
			Wrap(func() {})
		handler.ServeHTTP(rw, req)
		require.Equal(t, http.StatusCreated, rw.Result().StatusCode)

		handler = New().
			WithConstruct(nopConstructor).
			Before(func() *myerr {
				return nil
			}).
			Finally(func(rw http.ResponseWriter, res interface{}, err error) {
				require.NotNil(t, rw)
				require.NoError(t, err)
				rw.WriteHeader(http.StatusCreated)
			}).
			Wrap(func() {})
		handler.ServeHTTP(rw, req)
		require.Equal(t, http.StatusCreated, rw.Result().StatusCode)
	})
}
