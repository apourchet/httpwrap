package httpwrap

import (
	"fmt"
	"net/http"
	"reflect"
)

var _httpResponseWriterType = reflect.TypeOf(http.ResponseWriter(nil))
var _httpRequestType = reflect.TypeOf(&http.Request{})

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
	fmt.Println("constructing", t)
	if t.Kind() == reflect.Interface {
		return reflect.Zero(t), nil
	}
	obj := reflect.New(t).Elem()
	err := ctx.cons(ctx.w, ctx.req, obj)
	return obj, err
}
