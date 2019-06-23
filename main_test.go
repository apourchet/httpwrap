package httpwrap

import (
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMain(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test", nil)
		rw := httptest.NewRecorder()
		ctx := newRunCtx(rw, req, nopConstructor)

		main, err := newMain(func(in interface{}) error {
			require.Nil(t, in)
			return nil
		})
		require.NoError(t, err)
		res := main.run(ctx)
		require.Nil(t, res)
	})

	t.Run("with result", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test", nil)
		rw := httptest.NewRecorder()
		ctx := newRunCtx(rw, req, nopConstructor)

		type resp struct{ i int }

		main, err := newMain(func(in interface{}) (resp, error) {
			require.Nil(t, in)
			return resp{1}, nil
		})
		require.NoError(t, err)
		res := main.run(ctx)
		require.NotNil(t, res)
	})
}
