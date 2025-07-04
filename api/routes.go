package api

import (
	"net/http"
	"main/api/service"
)
func Route(){
	http.HandleFunc("/post",service.CreateEntry)

	http.HandleFunc("/",service.GetEntries)

	http.HandleFunc("/put/",service.UpdateEntries)

	http.HandleFunc("/delete/",service.DeleteEntries)
}