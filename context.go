package httpwrap

import (
	"net/http"
	"reflect"
)

type runctx struct {
	rw   http.ResponseWriter
	req  *http.Request
	cons RequestReader

	response    reflect.Value
	results     map[reflect.Type]param
	resultSlice []param
}

type param struct {
	t reflect.Type
	v reflect.Value
	i any
}

func newRunCtx(
	rw http.ResponseWriter,
	req *http.Request,
	cons RequestReader,
) *runctx {
	ctx := &runctx{
		req:         req,
		rw:          rw,
		cons:        cons,
		response:    reflect.Zero(reflect.TypeOf((*any)(nil)).Elem()),
		results:     map[reflect.Type]param{},
		resultSlice: []param{},
	}
	ctx.provide(req)
	ctx.provide(rw)
	return ctx
}

func (ctx *runctx) provide(i any) {
	if i == nil {
		return
	}
	p := param{
		t: reflect.TypeOf(i),
		v: reflect.ValueOf(i),
		i: i,
	}
	switch p.t.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
		if p.v.IsNil() {
			return
		}
	}
	ctx.results[p.t] = p
	ctx.resultSlice = append(ctx.resultSlice, p)
}

func (ctx *runctx) get(t reflect.Type) (val reflect.Value, found bool) {
	if isEmptyInterface(t) {
		if ctx.response.IsValid() {
			return ctx.response, true
		}
		return ctx.response, false
	}

	if t.Kind() != reflect.Interface {
		param, found := ctx.results[t]
		return param.v, found
	}

	for i := len(ctx.resultSlice) - 1; i >= 0; i-- {
		p := ctx.resultSlice[i]
		if p.t.Implements(t) {
			return p.v, true
		}
	}
	return val, false
}

func (ctx *runctx) construct(t reflect.Type) (reflect.Value, error) {
	if t.Kind() == reflect.Interface {
		return reflect.Zero(t), nil
	}

	// Make sure that we create a non-nil value if possible.
	obj := reflect.New(t)
	switch t.Kind() {
	case reflect.Ptr:
		obj.Elem().Set(reflect.New(t.Elem()))
	case reflect.Map:
		obj.Elem().Set(reflect.MakeMap(t))
	case reflect.Slice:
		obj.Elem().Set(reflect.MakeSlice(t, 0, 1))
	}

	err := ctx.cons(ctx.rw, ctx.req, obj.Interface())
	return obj.Elem(), err
}

func (ctx *runctx) generate(types []reflect.Type) ([]reflect.Value, error) {
	values := make([]reflect.Value, len(types))
	for i, t := range types {
		if val, found := ctx.get(t); found {
			values[i] = val
			continue
		}

		val, err := ctx.construct(t)
		if err != nil {
			ctx.provide(err)
			return nil, err
		}
		values[i] = val
	}
	return values, nil
}
