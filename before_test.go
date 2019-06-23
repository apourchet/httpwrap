package httpwrap

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBefore(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		before, err := newBefore(func(req *http.Request, rw http.ResponseWriter, in struct{}) error {
			return nil
		})
		require.NoError(t, err)

		req := httptest.NewRequest("GET", "/test", nil)
		rw := httptest.NewRecorder()
		ctx := newRunCtx(rw, req, nopConstructor)

		err = before.run(ctx)
		require.NoError(t, err)
	})
}
