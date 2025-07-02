package main

import (
	"log"
	"main/dbConn"
	"main/routes"
	"net/http"
	_ "github.com/lib/pq"
)

func main(){
	db:= dbconn.ConnToDb()
	routes.DB = db
	defer db.Close()

	http.HandleFunc("/post",routes.CreateEntry)

	http.HandleFunc("/",routes.GetEntries)

	http.HandleFunc("/put/",routes.UpdateEntries)

	http.HandleFunc("/delete/",routes.DeleteEntries)
	
	log.Fatal(http.ListenAndServe(":8000",nil))
}
