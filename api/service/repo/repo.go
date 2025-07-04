package repo

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Todo struct{
	Title string `json:"title"`
	Description string `json:"description"`
}
var db *sql.DB
func ConnToDb()*sql.DB{
	var err error
	err = godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}
	connStr:= os.Getenv("CONNECTION_STRING")
	
	db,err = sql.Open("postgres",connStr)
	checkErr(err)
	// createTable(db)
	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
	return db
}


// This function will insert into the DB table
func InsertIntoTable(todo Todo) int{
	query := `INSERT INTO todos (title,description)
	VALUES ($1,$2) RETURNING id`
	var pk int
	// fmt.Println(todo.Name)
	err :=db.QueryRow(query, todo.Title,todo.Description).Scan(&pk)
	checkErr(err)
	return pk
}

// This function will update the table query
func UpdateTable(id int, todo Todo)int64{
	query := `UPDATE todos SET title=$1, description=$2 WHERE id=$3`
	res,err := db.Exec(query,todo.Title,todo.Description,id)
	checkErr(err)
	rowsAffected,err:= res.RowsAffected()
	checkErr(err)
	return rowsAffected
}

// This function will delete the query
func DeleteEntry(id int)int64{
	query := `DELETE FROM todos WHERE id=$1`
	res,err := db.Exec(query,id)
	checkErr(err)
	rowsAffected,err:= res.RowsAffected()
	checkErr(err)
	return rowsAffected
}

// This function will print all the queries
func PrintTable()[]Todo{
	query:= "SELECT id,title,description FROM todos ORDER BY id"
	rows, err := db.Query(query)
	var list []Todo
	checkErr(err)
	for rows.Next() {
		var id int
		var title string
		var description string
		if err := rows.Scan(&id, &title,&description); err != nil {
			log.Fatal(err)
		}
		fmt.Println(id)
		todo:= Todo{Title:title,Description:description}
		list = append(list,todo)
	}
	return list
}


func checkErr(err error){
	if err!=nil{
		log.Fatal(err)
	}
}