package routes

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"main/sqlqueries"
	"net/http"
	"strconv"
	"strings"
	_ "github.com/lib/pq"
)

var DB *sql.DB

func CreateEntry(w http.ResponseWriter,r *http.Request){
	body,err:= io.ReadAll(r.Body)
	checkErr(err)
	todo:=sqlqueries.Todo{}
	err = json.NewDecoder(bytes.NewBuffer(body)).Decode(&todo)
	checkErr(err)
	pk:=sqlqueries.InsertIntoTable(DB,todo)
	fmt.Println(pk)
}

// This function is used to create a DB table
func createTable(db *sql.DB){
	query := `CREATE TABLE IF NOT EXISTS todos (
		id SERIAL PRIMARY KEY,
		title VARCHAR(255) NOT NULL,
		description VARCHAR(255)
	)`
	_,err :=db.Exec(query) 
	checkErr(err)
	fmt.Println("Table successfully created")
}

// This function will fetch all the entries from the table
func GetEntries(w http.ResponseWriter,r *http.Request){
	list:=sqlqueries.PrintTable(DB)
	w.Header().Set("Content-Type","application/json")
	err := json.NewEncoder(w).Encode(list)
	if err!=nil{
		fmt.Println("error")
		return
	}
}

// This function will update the row
func UpdateEntries(w http.ResponseWriter,r *http.Request){
	body,_ := io.ReadAll(r.Body)
	todo:= sqlqueries.Todo{}
	err := json.NewDecoder(bytes.NewReader(body)).Decode(&todo)
	if err!=nil{
		fmt.Println("error")
		return
	}
	path:= r.URL.Path
	parts:= strings.Split(path,"/")
	if len(parts)<3{
		fmt.Println("error! id not given")
		return
	}
	id,_:= strconv.Atoi(parts[2])
	sqlqueries.UpdateTable(DB,id,todo)
}

// This function will delete the row 
func DeleteEntries(w http.ResponseWriter,r *http.Request){
	path:= r.URL.Path
	parts:= strings.Split(path,"/")
	if len(parts)<3{
		fmt.Println("error! id not given")
		return
	}
	id,_:= strconv.Atoi(parts[2])
	sqlqueries.DeleteEntry(DB,id)
}

func checkErr(err error){
	if err!=nil{
		log.Fatal(err)
	}
}