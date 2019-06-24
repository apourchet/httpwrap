package httpwrap

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func nopConstructor(http.ResponseWriter, *http.Request, interface{}) error { return nil }

func jsonBodyConstructor(_ http.ResponseWriter, req *http.Request, obj interface{}) error {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, obj)
	req.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	return err
}

func failedConstructor(http.ResponseWriter, *http.Request, interface{}) error {
	return fmt.Errorf("error")
}

type myerr struct{}

func (e myerr) Error() string { return "error" }

var _ error = myerr{}

func TestContext(t *testing.T) {
	t.Run("default http types", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test", nil)
		rw := httptest.NewRecorder()
		ctx := newRunCtx(rw, req, nopConstructor)

		types, _ := typesOf(func(http.ResponseWriter, *http.Request) {})
		_, found := ctx.get(types[0])
		require.True(t, found)

		_, found = ctx.get(types[1])
		require.True(t, found)

		vals, err := ctx.generate(types)
		require.NoError(t, err)
		require.Len(t, vals, 2)
	})

	t.Run("provide error", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test", nil)
		rw := httptest.NewRecorder()
		ctx := newRunCtx(rw, req, nopConstructor)

		err := fmt.Errorf("error")
		ctx.provide(err)

		fn := func(err error) {}
		types, _ := typesOf(fn)

		vals, err := ctx.generate(types)
		require.NoError(t, err)
		require.Len(t, vals, 1)
		require.False(t, vals[0].IsNil())
	})

	t.Run("provide special error", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test", nil)
		rw := httptest.NewRecorder()
		ctx := newRunCtx(rw, req, nopConstructor)

		err := myerr{}
		ctx.provide(err)

		fn := func(err error) {}
		types, _ := typesOf(fn)

		vals, err1 := ctx.generate(types)
		require.NoError(t, err1)
		require.Len(t, vals, 1)
		require.NotNil(t, vals[0].Interface())
	})

	t.Run("provide nil value that satisfies interface", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test", nil)
		rw := httptest.NewRecorder()
		ctx := newRunCtx(rw, req, nopConstructor)

		var err *myerr
		ctx.provide(err)

		fn := func(err error) {}
		types, _ := typesOf(fn)

		vals, err1 := ctx.generate(types)
		require.NoError(t, err1)
		require.Len(t, vals, 1)
		require.Nil(t, vals[0].Interface())
	})
}
