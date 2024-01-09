package httpwrap

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAfter(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test", nil)
		rw := httptest.NewRecorder()
		ctx := newRunCtx(rw, req, nopConstructor)

		after, err := newAfter(func(w http.ResponseWriter, res any) {
			require.NotNil(t, w)
			w.WriteHeader(http.StatusOK)
		})
		require.NoError(t, err)
		after.run(ctx)

		require.Equal(t, http.StatusOK, rw.Result().StatusCode)
	})

	t.Run("with response", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test", nil)
		rw := httptest.NewRecorder()
		ctx := newRunCtx(rw, req, nopConstructor)
		ctx.response = reflect.ValueOf(1)

		after, err := newAfter(func(w http.ResponseWriter, res any) {
			require.NotNil(t, w)
			require.Equal(t, res, 1)
		})
		require.NoError(t, err)
		after.run(ctx)
	})
}
