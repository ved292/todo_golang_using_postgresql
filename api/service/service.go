package service

import (
	"bytes"
	"encoding/json"
	"io"
	"main/api/service/repo"
	"net/http"
	"strconv"
	"strings"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)


func Create(w http.ResponseWriter,r *http.Request){
	body,err:= io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	todo:=repo.Todo{}
	err = json.NewDecoder(bytes.NewBuffer(body)).Decode(&todo)
	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}
	_,err=repo.Create(todo)
	if err!=nil{
		http.Error(w,"Unable to query from the database",http.StatusInternalServerError)
		return
	}
	// fmt.Println(pk)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Entry Posted"))
}


// This function will fetch all the entries from the table
func Get(w http.ResponseWriter,r *http.Request){
	list,err:=repo.Get()
	if err!=nil{
		if strings.Contains(err.Error(),"failed to execute the query"){
			http.Error(w,"Failed the execute the query",http.StatusInternalServerError)
		}else if strings.Contains(err.Error(),"failed to get the rows"){
			http.Error(w,"Failed to get the rows",http.StatusInternalServerError)
		}
		return
	}
	w.Header().Set("Content-Type","application/json")
	err = json.NewEncoder(w).Encode(list)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// This function will update the row
func Update(w http.ResponseWriter,r *http.Request){
	body,err:= io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	todo:= repo.Todo{}
	err = json.NewDecoder(bytes.NewReader(body)).Decode(&todo)
	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}
	vars := mux.Vars(r)
	idStr := vars["id"]
	id,err:= strconv.Atoi(idStr)
	if err!=nil{
		http.Error(w,"The given id is not an integer",http.StatusBadRequest)
		return
	}
	rowsAffected,err:=repo.Update(id,todo)
	if err!=nil{
		if strings.Contains(err.Error(), "failed to execute update query") {
			http.Error(w, "DB query execution failed", http.StatusInternalServerError)
			
		}else if strings.Contains(err.Error(), "failed to get the affected rows"){
			http.Error(w, "Failed to get the affected rows", http.StatusInternalServerError)
		}
		return
	}
	if rowsAffected == 0 {
		http.Error(w, "Entry not found",404)
	} else {
		w.Write([]byte("Entry Updated"))
	}

}

// This function will delete the row 
func Delete(w http.ResponseWriter,r *http.Request){
	vars := mux.Vars(r)
	idStr := vars["id"]
	id,err:= strconv.Atoi(idStr)
	if err!=nil{
		http.Error(w,"The given id is not an integer",http.StatusBadRequest)
		return
	}
	rowsAffected,err:=repo.Delete(id)
	if err!=nil{
		if strings.Contains(err.Error(),"failed to execute update query"){
			http.Error(w,"Failed to execute delete query",http.StatusInternalServerError)
		}else if strings.Contains(err.Error(),"failed to get rows affected"){
			http.Error(w,"Failed to get rows affected",http.StatusInternalServerError)
		}
		return
	}
	if rowsAffected == 0 {
		http.Error(w, "Entry not found",404)
	} else {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Entry Deleted"))
	}
	
}

// func checkErr(err error){
// 	if err!=nil{
// 		log.Fatal(err)
// 	}
// }