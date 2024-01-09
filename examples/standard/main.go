package main

import (
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

func (pet Pet) IsInCategories(categories []int) bool {
	for _, c := range categories {
		if pet.Category == c {
			return true
		}
	}
	return false
}

var ErrBadAPICreds = httpwrap.NewHTTPError(http.StatusUnauthorized, "bad API credentials")
var ErrPetConflict = httpwrap.NewHTTPError(http.StatusConflict, "duplicate pet")
var ErrPetNotFound = httpwrap.NewHTTPError(http.StatusNotFound, "pet not found")

// ***** Middleware Definitions *****
// checkAPICreds checks the api credentials passed into the request.
func checkAPICreds(creds APICredentials) error {
	if creds.Key == "my-secret-key" {
		return nil
	}
	return ErrBadAPICreds
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

type FilterPetParams struct {
	Categories *[]int `http:"query=categories"`
	HasPhotos  *bool  `http:"query=hasPhotos"`
}

// FilterPets returns a list of pets that match the parameters given.
func (h *PetStoreHandler) FilterPets(params FilterPetParams) []Pet {
	res := []Pet{}
	for _, pet := range h.pets {
		if params.HasPhotos != nil && len(pet.PhotoURLs) == 0 {
			continue
		} else if params.Categories != nil && !pet.IsInCategories(*params.Categories) {
			continue
		}
		res = append(res, *pet)
	}
	return res
}

func (h *PetStoreHandler) ClearStore() error {
	h.pets = map[string]*Pet{}
	return nil
}

func main() {
	router := mux.NewRouter()

	handler := &PetStoreHandler{pets: map[string]*Pet{}}
	wrapper := httpwrap.NewStandardWrapper().Before(checkAPICreds)

	router.Handle("/pets", wrapper.Wrap(handler.AddPet)).Methods("POST")
	router.Handle("/pets", wrapper.Wrap(handler.GetPets)).Methods("GET")
	router.Handle("/pets/filtered", wrapper.Wrap(handler.FilterPets)).Methods("GET")
	router.Handle("/pets/{name}", wrapper.Wrap(handler.GetPetByName)).Methods("GET")
	router.Handle("/pets/{name}", wrapper.Wrap(handler.UpdatePet)).Methods("PUT")

	router.Handle("/clear", wrapper.Wrap(handler.ClearStore)).Methods("POST")

	http.Handle("/", router)
	log.Fatal(http.ListenAndServe(":3000", router))
}
