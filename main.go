package main

import (
	"fmt"
	"log"
	"main/api"
	"main/api/service/repo"
	"net/http"
	"strings"

	_ "github.com/lib/pq"
)

func main(){
	err:=repo.ConnToDb()
	if err!=nil{
		if strings.Contains(err.Error(),"failed to connect to Database"){
			fmt.Println("Error connecting with the database")
			return
		}
	}
	defer repo.CloseDB()
	
	api.Route()
	
	log.Fatal(http.ListenAndServe(":8000",nil))
}
