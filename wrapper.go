package httpwrap

import (
	"net/http"
	"reflect"
)

type Wrapper struct {
	befores []beforeFn
	main    mainFn
	afters  []afterFn

	construct Constructor
}

type Constructor func(http.ResponseWriter, *http.Request, interface{}) error

func New() *Wrapper {
	return &Wrapper{
		construct: func(http.ResponseWriter, *http.Request, interface{}) error { return nil },
	}
}

func (w Wrapper) WithConstruct(cons Constructor) *Wrapper {
	w.construct = cons
	return &w
}

func (w Wrapper) Before(fns ...interface{}) *Wrapper {
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
	return &w
}

func (w Wrapper) Wrap(fn interface{}) *Wrapper {
	main, err := newMain(fn)
	if err != nil {
		panic(err)
	}
	w.main = main
	return &w
}

func (w Wrapper) After(fns ...interface{}) *Wrapper {
	afters := make([]afterFn, len(w.afters)+len(fns))
	copy(afters, w.afters)
	for i, after := range fns {
		helper, err := newAfter(after)
		if err != nil {
			panic(err)
		}
		afters[i+len(w.afters)] = helper
	}
	w.afters = afters
	return &w
}

func (w Wrapper) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	ctx := newRunCtx(rw, req, w.construct)
	err := w.serveBefores(ctx)
	if err == nil {
		ctx.response = reflect.ValueOf(w.main.run(ctx))
	}
	w.serveAfters(ctx)
}

func (w Wrapper) serveBefores(ctx *runctx) error {
	for _, before := range w.befores {
		if err := before.run(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (w Wrapper) serveAfters(ctx *runctx) {
	for _, after := range w.afters {
		after.run(ctx)
	}
}
