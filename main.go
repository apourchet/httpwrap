package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// To test.
type hand struct {
	i int
}

type r1 struct {
	Auth string `json:"auth"`
}

func (h hand) mw(auth r1) Result {
	if auth.Auth != "" {
		return Result{
			Body: map[string]interface{}{
				"subject": auth.Auth,
				"roles":   []string{"user", "staff"},
			},
		}
	}
	return Result{Done: true, StatusCode: 401, ErrExternal: fmt.Errorf("unauthorized")}
}

func (h hand) h1(ctx Context, req struct {
	Subject string
	Roles   []string
}) (res Result) {
	if req.Subject == "google|123" {
		res.StatusCode = http.StatusOK
	} else {
		res.StatusCode = http.StatusForbidden
	}
	return
}

func (h hand) h2() (res Result) {
	res.StatusCode = http.StatusOK
	return
}

func h3() (res Result) {
	res.StatusCode = http.StatusOK
	return
}

func h4(ctx Context) {
}

func main() {
	handler := hand{100}

	r := mux.NewRouter()
	r.Handle("/h1", Wrap(handler.mw, handler.h1))
	r.Handle("/h2", Wrap(handler.h2))
	r.Handle("/h3", Wrap(h3))
	r.Handle("/h3", Wrap(h4))

	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8081", r))
}
