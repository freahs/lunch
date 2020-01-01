package main

import (
	"encoding/json"
	"github.com/freahs/lunch-server/data"
	"github.com/freahs/lunch-server/loaders"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

// HTTPError is an error with a status
type HTTPError struct {
	status int
	error
}

// MarshalJSON implements the Marshaller interface
func (e HTTPError) MarshalJSON() ([]byte, error) {
	type E struct {
		S int    `json:"status"`
		M string `json:"message"`
	}
	return json.Marshal(E{e.status, e.Error()})
}

// Status returns the status code of the error
func (e HTTPError) Status() int {
	return e.status
}

// NewHTTPError returns a new HTTPError
func NewHTTPError(status int, err error) HTTPError {
	return HTTPError{status, err}
}

func loadStore() *data.Store {
	store, err := data.LoadStore("/home/fredrik/tmp/teststore.json")
	if err != nil {
		log.Fatal(err)
	}
	return store
}

func main() {

	store := data.NewStore()
	prime := loaders.NewPrime()
	err := prime.Load(store)
	if err != nil {
		log.Fatal(err)
	}
	Router := mux.NewRouter().StrictSlash(true)
	APIRouter := Router.PathPrefix("/api/v1").Subrouter()
	api := NewAPIServer(store, APIRouter)

	log.Fatal(http.ListenAndServe(":8080", api))
}
