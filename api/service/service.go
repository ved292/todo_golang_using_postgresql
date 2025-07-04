package service
import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"main/api/service/repo"
	"net/http"
	"strconv"
	"strings"
	_ "github.com/lib/pq"
)

var DB *sql.DB

func CreateEntry(w http.ResponseWriter,r *http.Request){
	body,err:= io.ReadAll(r.Body)
	checkErr(err)
	todo:=repo.Todo{}
	err = json.NewDecoder(bytes.NewBuffer(body)).Decode(&todo)
	checkErr(err)
	repo.InsertIntoTable(todo)
	// fmt.Println(pk)
	w.Write([]byte("Entry Posted"))
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
	list:=repo.PrintTable()
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
	todo:= repo.Todo{}
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
	rowsAffected:=repo.UpdateTable(id,todo)
	if rowsAffected == 0 {
		w.Write([]byte("Entry not found"))
	} else {
		w.Write([]byte("Entry Updated"))
	}
	
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
	rowsAffected:=repo.DeleteEntry(id)
	if rowsAffected == 0 {
		w.Write([]byte("Entry not found"))
	} else {
		w.Write([]byte("Entry Deleted"))
	}
	
}

func checkErr(err error){
	if err!=nil{
		log.Fatal(err)
	}
}