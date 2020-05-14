package httpwrap

import (
	"net/http"
	"reflect"
)

// Wrapper implements the http.Handler interface, wrapping the handlers
// that are passed in.
type Wrapper struct {
	befores   []beforeFn
	after     *afterFn
	construct Constructor
}

// New creates a new Wrapper object.
func New() Wrapper {
	return Wrapper{
		construct: EmptyConstructor,
	}
}

// WithConstruct returns a new wrapper with the given Constructor function.
func (w Wrapper) WithConstruct(cons Constructor) Wrapper {
	w.construct = cons
	return w
}

// Before adds a new function that will execute before the main handler. The chain
// of befores will end if a before returns a non-nil error value.
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

// Finally sets the last function that will execute during a request.
func (w Wrapper) Finally(fn interface{}) Wrapper {
	after, err := newAfter(fn)
	if err != nil {
		panic(err)
	}
	w.after = &after
	return w
}

// Wrap sets the main handling function to process requests. This Wrap function must
// be called to get an `http.Handler` type.
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

// Handler is a Wrapper that implements `http.Handler`.
type Handler struct {
	Wrapper
	main mainFn
}

// Before adds the before functions to the underlying Wrapper.
func (h Handler) Before(fns ...interface{}) Handler {
	h.Wrapper = h.Wrapper.Before(fns...)
	return h
}

// Finally sets the `finally` function of the underlying Wrapper.
func (h Handler) Finally(fn interface{}) Handler {
	h.Wrapper = h.Wrapper.Finally(fn)
	return h
}

// ServeHTTP implements `http.Handler`.
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
