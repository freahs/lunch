package main

import (
	"github.com/freahs/lunch"
	"github.com/freahs/lunch/loaders"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

func main() {

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	store := lunch.NewStore()
	prime := loaders.NewPrime()
	err := prime.Load(store)
	if err != nil {
		log.Fatal(err)
	}
	Router := mux.NewRouter().StrictSlash(true)
	APIRouter := Router.PathPrefix("/api/v1").Subrouter()
	api := NewAPIServer(store, APIRouter)
	log.Printf("starting api server @ port %v", port)
	log.Fatal(http.ListenAndServe(":"+port, api))
}
