package main

import (
	"encoding/json"
	"github.com/freahs/lunch-server/data"
	"github.com/freahs/lunch-server/loaders"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"path/filepath"
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

func loadStore(filename string) *data.Store {
	dir := func() string {
		if d, err := os.UserConfigDir(); err == nil {
			d = filepath.Join(d, "lunch-server")
			if err = os.MkdirAll(d, 0755); err == nil {
				return d
			}
		}
		return os.TempDir()
	}()
	path := filepath.Join(dir, filename)
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return data.NewStore()
	}
	if info.IsDir() {
		log.Fatalf("%v is a directory", path)
	}
	store, err := data.LoadStore(path)
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
