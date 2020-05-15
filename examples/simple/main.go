package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/apourchet/httpwrap"
	"github.com/gorilla/mux"
)

// ***** Type Definitions *****
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

var ErrBadAPICreds = fmt.Errorf("bad API credentials")
var ErrPetConflict = fmt.Errorf("duplicate pet")
var ErrPetNotFound = fmt.Errorf("pet not found")

// ***** Middleware Definitions *****
type Middlewares struct{}

// checkAPICreds checks the api credentials passed into the request.
func (mw *Middlewares) checkAPICreds(creds APICredentials) error {
	if creds.Key == "my-secret-key" {
		return nil
	}
	return ErrBadAPICreds
}

// sendResponse writes out the response to the client given the output
// of the handler.
func (mw *Middlewares) sendResponse(w http.ResponseWriter, res interface{}, err error) {
	switch err {
	case ErrBadAPICreds:
		w.WriteHeader(http.StatusUnauthorized)
	case ErrPetConflict:
		w.WriteHeader(http.StatusConflict)
	case ErrPetNotFound:
		w.WriteHeader(http.StatusNotFound)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}

	if err != nil {
		if _, err := w.Write([]byte(err.Error())); err != nil {
			log.Println("error writing response:", err)
		}
	} else {
		if err := json.NewEncoder(w).Encode(res); err != nil {
			log.Println("Error writing response:", err)
		}
	}
}

// ***** Handler Methods *****

// AddPet adds a new pet to the store.
func (h *PetStoreHandler) AddPet(pet Pet) error {
	if _, found := h.pets[pet.Name]; found {
		return ErrPetConflict
	}
	h.pets[pet.Name] = &pet
	return nil
}

// GetPets returns the list of pets in the store.
func (h *PetStoreHandler) GetPets() (res []Pet, err error) {
	res = make([]Pet, 0, len(h.pets))
	for _, pet := range h.pets {
		res = append(res, *pet)
	}
	return res, nil
}

type GetByNameParams struct {
	Name string `http:"segment=name"`
}

// GetPetByName returns a pet given its name.
func (h *PetStoreHandler) GetPetByName(params GetByNameParams) (pet *Pet, err error) {
	pet, found := h.pets[params.Name]
	if !found {
		return nil, ErrPetNotFound
	}
	return pet, nil
}

type UpdateParams struct {
	Name string `http:"segment=name"`

	Category  *int      `json:"category"`
	PhotoURLs *[]string `json:"photoUrls"`
}

// UpdatePet updates a pet given its name.
func (h *PetStoreHandler) UpdatePet(params UpdateParams) error {
	pet, found := h.pets[params.Name]
	if !found {
		return ErrPetNotFound
	}

	if params.Category != nil {
		pet.Category = *params.Category
	}
	if params.PhotoURLs != nil {
		pet.PhotoURLs = *params.PhotoURLs
	}
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

	r.Handle("/pets", wrapper.Wrap(handler.AddPet)).Methods("POST")
	r.Handle("/pets", wrapper.Wrap(handler.GetPets)).Methods("GET")
	r.Handle("/pets/{name}", wrapper.Wrap(handler.GetPetByName)).Methods("GET")
	r.Handle("/pets/{name}", wrapper.Wrap(handler.UpdatePet)).Methods("PUT")

	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":3000", r))
}
