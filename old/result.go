package httpwrap

import (
	"encoding/json"
	"net/http"
)

type Result struct {
	StatusCode int

	ErrInternal error
	ErrExternal error

	Body interface{}

	Done bool
}

func (r *Result) Finish() {
	r.Done = true
}

func (r Result) send(w http.ResponseWriter) error {
	if r.StatusCode != 0 {
		w.WriteHeader(r.StatusCode)
	}
	switch {
	case r.ErrExternal != nil:
		return json.NewEncoder(w).Encode(r.ErrExternal)
	case r.Body != nil:
		return json.NewEncoder(w).Encode(r.Body)
	}
	return nil
}
