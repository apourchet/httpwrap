package httpwrap

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func nopConstructor(http.ResponseWriter, *http.Request, interface{}) error { return nil }

func jsonBodyConstructor(_ http.ResponseWriter, req *http.Request, obj interface{}) error {
	return json.NewDecoder(req.Body).Decode(obj)
}

func TestContext(t *testing.T) {
	req := httptest.NewRequest("GET", "/test", nil)
	rw := httptest.NewRecorder()
	ctx := newRunCtx(rw, req, nopConstructor)

	_, found := ctx.results[reflect.TypeOf(http.ResponseWriter(nil))]
	require.True(t, found)

	_, found = ctx.results[reflect.TypeOf(&http.Request{})]
	require.True(t, found)
}
