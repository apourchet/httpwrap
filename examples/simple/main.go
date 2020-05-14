package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/apourchet/httpwrap"
	"github.com/gorilla/mux"
)

// Type Definitions
type APICredentials struct {
	Key string `http:"header=X-PETSTORE-KEY"`
}

type PetStoreHandler struct {
	pets map[string]*Pet
}

type Pet struct {
	Name      string   `json:"name"`
	Category  int      `json:"category"`
	PhotoURLs []string `json:"photoUrls"`
}

type PetParams struct {
	Name      *string   `json:"name"`
	Category  *int      `json:"category"`
	PhotoURLs *[]string `json:"photoUrls"`
}

var ErrBadAPICreds = fmt.Errorf("bad API credentials")

// Middleware definition
type Middlewares struct{}

// checkAPICreds checks the api credentials passed into the request.
func (mw *Middlewares) checkAPICreds(creds APICredentials) error {
	if creds.Key == "my-secret-key" {
		return nil
	}
	return ErrBadAPICreds
}

func (mw *Middlewares) sendResponse(w http.ResponseWriter, res interface{}, err error) {
	if err == ErrBadAPICreds {
		w.WriteHeader(http.StatusUnauthorized)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			log.Println("error writing response:", err)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		log.Println("Error writing response:", err)
	}
}

// Handler Methods
func (h *PetStoreHandler) AddPet(pet Pet) error {
	// TODO
	return nil
}

func (h *PetStoreHandler) GetPets() (res []Pet, err error) {
	// TODO
	return res, nil
}

func (h *PetStoreHandler) UpdatePet(params PetParams) error {
	// TODO
	return nil
}

func main() {
	r := mux.NewRouter()

	handler := &PetStoreHandler{pets: map[string]*Pet{}}
	mw := &Middlewares{}

	wrapper := httpwrap.New().
		WithConstruct(httpwrap.StandardConstructor()).
		Before(mw.checkAPICreds).
		Finally(mw.sendResponse)

	r.Handle("/pets", wrapper.Wrap(handler.GetPets)).Methods("GET")

	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":3000", r))
}
