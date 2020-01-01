package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/freahs/lunch-server/data"
	"github.com/gorilla/mux"
)

func UnmarshalRequestBody(r *http.Request, i interface{}) error {

	bts, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bts, i)
	if err != nil {
		return err
	}
	return nil
}

func WriteResponse(w http.ResponseWriter, status int, data interface{}) {
	if err, ok := data.(error); ok {
		w.WriteHeader(status)
		w.Write([]byte(err.Error()))
	} else if bts, err := json.Marshal(data); err != nil {
		WriteResponse(w, http.StatusInternalServerError, err)
	} else {
		w.WriteHeader(status)
		w.Write(bts)
	}
}

type APIServer struct {
	store  *data.Store
	router *mux.Router
}

func NewAPIServer(store *data.Store, router *mux.Router) APIServer {
	api := APIServer{store, router}
	api.router.HandleFunc("/menu", api.createMenu).Methods("POST")
	api.router.HandleFunc("/menu", api.allMenus).Methods("GET")
	api.router.HandleFunc("/menu/{name}", api.getRestaurant).Methods("GET")
	return api
}

func (api APIServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	api.router.ServeHTTP(w, r)
}

func (api APIServer) filterMenus(r *http.Request, store *data.Store) *data.Store {
	query := r.URL.Query()
	filters := []data.Filter{data.FilterLt, data.FilterLe, data.FilterEq, data.FilterGe, data.FilterGt}
	for i, filter := range []string{"lt", "le", "eq", "ge", "gt"} {
		val := query.Get(filter)
		if val == "" || len(val) != 8 {
			continue
		}
		year, err := strconv.Atoi(val[:4])
		if err != nil {
			continue
		}
		month, err := strconv.Atoi(val[4:6])
		if err != nil {
			continue
		}
		day, err := strconv.Atoi(val[6:8])
		if err != nil {
			continue
		}
		store = store.FilterDate(filters[i], year, month, day)
	}
	return store
}

func (api APIServer) createMenu(w http.ResponseWriter, r *http.Request) {
	var menu data.Menu
	if err := UnmarshalRequestBody(r, &menu); err != nil {
		WriteResponse(w, http.StatusUnprocessableEntity, err)
		return
	}
	api.store.AddMenu(menu)
	WriteResponse(w, http.StatusOK, menu)
}

func (api APIServer) allMenus(w http.ResponseWriter, r *http.Request) {
	WriteResponse(w, http.StatusOK, api.filterMenus(r, api.store))
}

func (api APIServer) getRestaurant(w http.ResponseWriter, r *http.Request) {
	if name, ok := mux.Vars(r)["name"]; !ok {
		WriteResponse(w, http.StatusInternalServerError, fmt.Errorf("no such variable 'name'"))
	} else {
		store := api.store.FilterName(name)
		store = api.filterMenus(r, store)
		WriteResponse(w, http.StatusOK, store)
	}
}
