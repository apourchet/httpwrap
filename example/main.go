package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/apourchet/httpwrap"
	"github.com/gorilla/mux"
)

type hand struct {
	i int
}

type AuthReq struct {
	Secret string
}

type Description struct {
	Subject string
	Roles   []string
}

func (h hand) finish(w http.ResponseWriter, res interface{}, err error) {
	fmt.Println("finish", res, err)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			log.Println("Error writing response:", err)
		}
		return
	}
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		log.Println("Error writing response:", err)
	}
}

func (h hand) describe(req *AuthReq) (*Description, error) {
	fmt.Println("describe", req)
	if req.Secret == "staff" {
		return &Description{
			Subject: "antoine",
			Roles:   []string{"user", "staff"},
		}, nil
	} else if req.Secret == "default" {
		return &Description{
			Subject: "default",
			Roles:   []string{"user"},
		}, nil
	}
	return nil, fmt.Errorf("unauthorized")
}

func (h hand) ensureStaff(desc *Description) error {
	fmt.Println("ensurestaff", desc)
	if len(desc.Roles) != 2 || desc.Roles[1] != "staff" {
		return fmt.Errorf("forbidden")
	}
	return nil
}

type H1Req struct {
	Stuff int
}

type H1Res struct {
	ID int
}

func (h hand) h1(desc Description, req H1Req) (*H1Res, error) {
	fmt.Println("h1", desc, req)
	if desc.Subject == "default" {
		return nil, fmt.Errorf("must be logged in")
	}
	fmt.Println("h1 => ", req.Stuff)
	return &H1Res{123}, nil
}

func construct(rw http.ResponseWriter, req *http.Request, obj interface{}) error {
	return json.NewDecoder(req.Body).Decode(obj)
}

func main() {
	handler := hand{100}

	wrapper := httpwrap.New().
		WithConstruct(construct).
		Before(handler.describe)

	r := mux.NewRouter()
	r.Handle("/h1", wrapper.
		Wrap(handler.h1).
		Finally(handler.finish))
	r.Handle("/h2", wrapper.
		Before(handler.ensureStaff).
		Wrap(handler.h1).
		Finally(handler.finish))

	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8081", r))
}
