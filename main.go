package main

import (
	"log"
	"main/api"
	"net/http"
	"main/api/service/repo"
	_ "github.com/lib/pq"
)

func main(){
	db:= repo.ConnToDb()
	defer db.Close()
	api.Route()
	
	log.Fatal(http.ListenAndServe(":8000",nil))
}
