package httpwrap

import (
	"net/http"
	"reflect"
)

type runctx struct {
	w    http.ResponseWriter
	req  *http.Request
	cons Constructor

	response reflect.Value
	results  map[reflect.Type]reflect.Value
}

func newRunCtx(
	w http.ResponseWriter,
	req *http.Request,
	cons Constructor,
) *runctx {
	ctx := &runctx{
		req:      req,
		w:        w,
		cons:     cons,
		response: reflect.ValueOf(nil),
		results: map[reflect.Type]reflect.Value{
			reflect.TypeOf(w):   reflect.ValueOf(w),
			reflect.TypeOf(req): reflect.ValueOf(req),
		},
	}
	return ctx
}

func (ctx *runctx) provide(t reflect.Type, val reflect.Value) {
	ctx.results[t] = val
}

func (ctx *runctx) construct(t reflect.Type) (reflect.Value, error) {
	// TODO: Use constructor to build reflect.Value.
	return reflect.New(t), nil
}
