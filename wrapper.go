package httpwrap

import (
	"net/http"
	"reflect"
)

// Wrapper implements the http.Handler interface, wrapping the handlers
// that are passed in
type Wrapper struct {
	befores   []beforeFn
	after     *afterFn
	construct Constructor
}

type Constructor func(http.ResponseWriter, *http.Request, interface{}) error

func New() Wrapper {
	return Wrapper{
		construct: func(http.ResponseWriter, *http.Request, interface{}) error { return nil },
	}
}

func (w Wrapper) WithConstruct(cons Constructor) Wrapper {
	w.construct = cons
	return w
}

func (w Wrapper) Before(fns ...interface{}) Wrapper {
	befores := make([]beforeFn, len(w.befores)+len(fns))
	copy(befores, w.befores)
	for i, before := range fns {
		helper, err := newBefore(before)
		if err != nil {
			panic(err)
		}
		befores[i+len(w.befores)] = helper
	}
	w.befores = befores
	return w
}

func (w Wrapper) Finally(fn interface{}) Wrapper {
	after, err := newAfter(fn)
	if err != nil {
		panic(err)
	}
	w.after = &after
	return w
}

func (w Wrapper) Wrap(fn interface{}) Handler {
	main, err := newMain(fn)
	if err != nil {
		panic(err)
	}
	return Handler{
		Wrapper: w,
		main:    main,
	}
}

type Handler struct {
	Wrapper
	main mainFn
}

func (h Handler) Before(fns ...interface{}) Handler {
	h.Wrapper = h.Wrapper.Before(fns...)
	return h
}

func (h Handler) Finally(fn interface{}) Handler {
	h.Wrapper = h.Wrapper.Finally(fn)
	return h
}

func (h Handler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	ctx := newRunCtx(rw, req, h.construct)
	err := h.serveBefores(ctx)
	if err == nil {
		ctx.response = reflect.ValueOf(h.main.run(ctx))
	}
	if h.after != nil {
		h.after.run(ctx)
	}
}

func (h Handler) serveBefores(ctx *runctx) error {
	for _, before := range h.befores {
		if err := before.run(ctx); err != nil {
			return err
		}
	}
	return nil
}
