package httpwrap

import (
	"errors"
	"net/http"
	"reflect"
)

// Wrapper implements the http.Handler interface, wrapping the handlers
// that are passed in.
type Wrapper struct {
	befores   []beforeFn
	after     *afterFn
	construct RequestReader
}

// New creates a new Wrapper object. This wrapper object will not interact in any way
// with the http request and response writer.
func New() Wrapper {
	return Wrapper{
		construct: emptyRequestReader,
	}
}

// WithRequestReader returns a new wrapper with the given RequestReader function.
func (w Wrapper) WithRequestReader(cons RequestReader) Wrapper {
	w.construct = cons
	return w
}

// Before adds a new function that will execute before the main handler. The chain
// of befores will end if a before returns a non-nil error value.
func (w Wrapper) Before(fns ...any) Wrapper {
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

// Finally sets the last function that will execute during a request. This function gets
// invoked with the response object and the possible error returned from the main
// endpoint function.
func (w Wrapper) Finally(fn any) Wrapper {
	after, err := newAfter(fn)
	if err != nil {
		panic(err)
	}
	w.after = &after
	return w
}

// Wrap sets the main handling function to process requests. This Wrap function must
// be called to get an `http.Handler` type.
func (w Wrapper) Wrap(fn any) http.Handler {
	main, err := newMain(fn)
	if err != nil {
		panic(err)
	}
	return wrappedHttpHandler{
		Wrapper: w,
		main:    main,
	}
}

// wrappedHttpHandler is a Wrapper that implements `http.Handler`.
type wrappedHttpHandler struct {
	Wrapper
	main mainFn
}

// ServeHTTP implements `http.Handler`.
func (h wrappedHttpHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	ctx := newRunCtx(rw, req, h.construct)
	err := h.serveBefores(ctx)
	if err == nil {
		ctx.response = reflect.ValueOf(h.main.run(ctx))
	}
	if h.after != nil {
		h.after.run(ctx)
	}
}

func (h wrappedHttpHandler) serveBefores(ctx *runctx) error {
	httpResponseType := reflect.TypeOf((*HTTPResponse)(nil)).Elem()
	for _, before := range h.befores {
		if err := before.run(ctx); err != nil {
			return err
		}
		possibleHttpResponses, err := ctx.generate([]reflect.Type{
			httpResponseType,
		})
		if err != nil {
			continue
		}

		for _, response := range possibleHttpResponses {
			if ok := response.Type().Implements(httpResponseType); ok {
				if !response.IsNil() {
					return errors.New("early middleware return")
				}
			}
		}
	}
	return nil
}
