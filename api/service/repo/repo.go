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
	Id int `json:"id,omitempty"`
	Title string `json:"title"`
	Description string `json:"description"`
}


const (
	host = "localhost"
	port = 5432
	user = "postgres"
	dbname = "todos"
)
var db *sql.DB
func ConnToDb()(error){
	var err error
	err = godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}
	// connStr:= os.Getenv("CONNECTION_STRING")
	password:= os.Getenv("PASSWORD")
	connStr:= fmt.Sprintf("host= %s port= %d user = %s password = %s dbname = %s sslmode=disable",host,port,user,password,dbname)
	db,err = sql.Open("postgres",connStr)
	if err!=nil{
		return fmt.Errorf("failed to connect to Database: %w", err)
	}
	// createTable(db)
	if err = db.Ping(); err != nil {
		return fmt.Errorf("failed to do connection with database: %w", err)
	}
	return nil
}

func CloseDB() {
	if db != nil {
		db.Close()
	}
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

// This function will insert into the DB table
func Create(todo Todo) (int,error){
	query := `INSERT INTO todos (title,description)
	VALUES ($1,$2) RETURNING id`
	var pk int
	// fmt.Println(todo.Name)
	err :=db.QueryRow(query, todo.Title,todo.Description).Scan(&pk)
	if err != nil {
		return 0, fmt.Errorf("failed to execute create query: %w", err)
	}
	return pk,nil
}

// This function will update the table query
func Update(id int, todo Todo)(int64,error){
	query := `UPDATE todos SET title=$1, description=$2 WHERE id=$3`
	res,err := db.Exec(query,todo.Title,todo.Description,id)
	if err != nil {
		return 0, fmt.Errorf("failed to execute update query: %w", err)
	}
	rowsAffected,err:= res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get the affected rows: %w", err)
	}
	return rowsAffected,nil
}



// This function will delete the query
func Delete(id int)(int64,error){
	query := `DELETE FROM todos WHERE id=$1`
	res,err := db.Exec(query,id)
	if err != nil {
		return 0, fmt.Errorf("failed to execute delete query: %w", err)
	}
	rowsAffected,err:= res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get affected rows: %w", err)
	}
	return rowsAffected,nil
}

// This function will print all the queries
func Get()([]Todo,error){
	query:= "SELECT id,title,description FROM todos ORDER BY id"
	rows, err := db.Query(query)
	var list []Todo
	if err!=nil{
		return list,fmt.Errorf("failed to execute the query: %w",err)
	}
	for rows.Next() {
		var id int
		var title string
		var description string
		if err := rows.Scan(&id, &title,&description); err != nil {
			return list,fmt.Errorf("failed to get the rows: %w",err)
		}
		// fmt.Println(id)
		todo:= Todo{Id:id,Title:title,Description:description}
		list = append(list,todo)
	}
	return list,nil
}


func checkErr(err error){
	if err!=nil{
		log.Fatal(err)
	}
}