package repo

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// The pointer helps to find the client send which fields
type Todo struct{
	Id *int `json:"id,omitempty"`
	Title *string `json:"title,omitempty"`
	Description *string `json:"description,omitempty"`
	IsDeleted *bool `json:"is_deleted,omitempty"`
}


const (
	host = "db"
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
	for range 5 {
		db, err = sql.Open("postgres", connStr)
		if err == nil && db.Ping() == nil {
			break
		}
		log.Println("Waiting for database to be ready...")
		time.Sleep(3 * time.Second)
	}
	createTable(db)
	//createIndex()
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

func createIndex(){
	query := `CREATE INDEX idx_active_todos ON todo(id, title, description) WHERE is_deleted = false`
	_,err := db.Exec(query)
	checkErr(err)
	fmt.Println("Index successfully created")
}

// This function is used to create a DB table
func createTable(db *sql.DB){
	query := `CREATE TABLE IF NOT EXISTS todo (
		id SERIAL PRIMARY KEY,
		title VARCHAR(255) NOT NULL,
		description VARCHAR(255),
		is_deleted BOOLEAN
	)`
	_,err :=db.Exec(query) 
	checkErr(err)
	fmt.Println("Table successfully created")
}


// This function will insert into the DB table
func Create(todo Todo) (int,error){
	query := `INSERT INTO todo (title,description,is_deleted)
	VALUES ($1,$2,false) RETURNING id`
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
	setClauses := []string{}
	args := []any{}
	argPosition := 1

	if todo.Title != nil {
		setClauses = append(setClauses, fmt.Sprintf("title=$%d", argPosition))
		args = append(args, todo.Title)
		argPosition++
	}
	if todo.Description != nil {
		setClauses = append(setClauses, fmt.Sprintf("description=$%d", argPosition))
		args = append(args, todo.Description)
		argPosition++
	}

	if todo.IsDeleted != nil{
		setClauses = append(setClauses, fmt.Sprintf("is_deleted=$%d", argPosition))
		args = append(args, todo.IsDeleted)
		argPosition++
	}

	if len(setClauses) == 0 {
		return 0, fmt.Errorf("no fields provided for update")
	}

	query := fmt.Sprintf("UPDATE todo SET %s WHERE id=$%d",strings.Join(setClauses, ", "), argPosition)

	args = append(args, id)

	res, err := db.Exec(query, args...)
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
	query := `UPDATE todo SET is_deleted = true WHERE id=$1`
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
	query:= "SELECT * FROM todo WHERE is_deleted = false ORDER BY id"
	rows, err := db.Query(query)
	var list []Todo
	if err!=nil{
		return list,fmt.Errorf("failed to execute the query: %w",err)
	}
	for rows.Next() {
		var id int
		var title string
		var description string
		var is_deleted bool
		if err := rows.Scan(&id, &title,&description,&is_deleted); err != nil {
			return list,fmt.Errorf("failed to get the rows: %w",err)
		}
		// fmt.Println(id)
		todo:= Todo{Id:&id,Title:&title,Description:&description,IsDeleted: &is_deleted}
		list = append(list,todo)
	}
	return list,nil
}


func checkErr(err error){
	if err!=nil{
		log.Fatal(err)
	}
}