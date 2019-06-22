package httpwrap

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// Wrapper is the struct that acts as a http.Handler.
type Wrapper struct {
	handlers []handler
	catch    func(error)
}

// Wrap wraps the handler functions into a http.Handler.
func Wrap(fn interface{}, fns ...interface{}) *Wrapper {
	if fn == nil {
		panic("cannot wrap nil function")
	}
	handlers := make([]handler, 1+len(fns))
	handlers[0] = newHandler(fn)
	for i, f := range fns {
		handlers[i+1] = newHandler(f)
	}
	return &Wrapper{
		handlers: handlers,
		catch: func(err error) {
			log.Println(err)
		},
	}
}

// WithCatch tells the wrapper which logging function it should use
// to log errors outside of the http request lifecycle.
func (wrapper *Wrapper) WithCatch(catch func(error)) *Wrapper {
	wrapper.catch = catch
	return wrapper
}

// ServeHTTP implements the http.Handler interface.
func (wrapper *Wrapper) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	ctx, err := newContext(w, req)
	if err != nil {
		err := Result{ErrExternal: err}.send(w)
		if err != nil && wrapper.catch != nil {
			wrapper.catch(err)
		}
		return
	}

	var lastResult *Result
	for _, h := range wrapper.handlers {
		result := h.handle(ctx)
		if result == nil {
			continue
		}

		ctx.Results = append(ctx.Results, *result)
		lastResult = result
		if result.Done {
			break
		}

		extra, err := json.Marshal(result.Body)
		if err != nil {
			err := fmt.Errorf("failed to marshal body: %v", err)
			ctx.Results = append(ctx.Results, Result{ErrInternal: err})
		}
		ctx.bodies = append(ctx.bodies, extra)
	}

	if lastResult == nil {
		return
	} else if err := lastResult.send(w); err != nil && wrapper.catch != nil {
		wrapper.catch(err)
	}
}
