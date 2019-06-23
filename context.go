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
		response: reflect.Zero(reflect.TypeOf((*interface{})(nil)).Elem()),
		results: map[reflect.Type]reflect.Value{
			_httpResponseWriterType: reflect.ValueOf(w),
			_httpRequestType:        reflect.ValueOf(req),
		},
	}
	return ctx
}

func (ctx *runctx) provide(t reflect.Type, val reflect.Value) {
	ctx.results[t] = val
}

func (ctx *runctx) get(t reflect.Type) (reflect.Value, bool) {
	if t.String() == "http.ResponseWriter" {
		return reflect.ValueOf(ctx.w), true
	}
	val, found := ctx.results[t]
	return val, found
}

func (ctx *runctx) construct(t reflect.Type) (reflect.Value, error) {
	if t.Kind() == reflect.Interface {
		return reflect.Zero(t), nil
	}
	obj := reflect.New(t).Elem()
	err := ctx.cons(ctx.w, ctx.req, obj)
	return obj, err
}
