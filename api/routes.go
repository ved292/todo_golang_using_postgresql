package api

import (
	"net/http"
	"main/api/service"
	"github.com/gorilla/mux"
)
func Route(){
	r:= mux.NewRouter()

	r.HandleFunc("/",service.Create).Methods("POST")
	
	r.HandleFunc("/",service.Get).Methods("GET")

	r.HandleFunc("/{id}",service.Update).Methods("PUT")

	r.HandleFunc("/{id}",service.Delete).Methods("DELETE")

	http.Handle("/",r)
}