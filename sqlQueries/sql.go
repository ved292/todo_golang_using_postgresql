package sqlqueries

import (
	"database/sql"
	"log"
	_ "github.com/lib/pq"
)

type Todo struct{
	Title string `json:"title"`
	Description string `json:"description"`
}

type TodoFetch struct{
	Id int `json:"id"`
	Title string `json:"title"`
	Description string `json:"description"`
}

// This function will insert into the DB table
func InsertIntoTable(db *sql.DB,todo Todo) int{
	query := `INSERT INTO todos (title,description)
	VALUES ($1,$2) RETURNING id`
	var pk int
	// fmt.Println(todo.Name)
	err :=db.QueryRow(query, todo.Title,todo.Description).Scan(&pk)
	checkErr(err)
	return pk
}

// This function will update the table query
func UpdateTable(db *sql.DB, id int, todo Todo){
	query := `UPDATE todos SET title=$1, description=$2 WHERE id=$3`
	_,err := db.Exec(query,todo.Title,todo.Description,id)
	checkErr(err)
}

// This function will delete the query
func DeleteEntry(db *sql.DB, id int){
	query := `DELETE FROM todos WHERE id=$1`
	_,err := db.Exec(query,id)
	checkErr(err)
}

// This function will print all the queries
func PrintTable(db *sql.DB)[]TodoFetch{
	query:= "SELECT id,title,description FROM todos ORDER BY id"
	rows, err := db.Query(query)
	var list []TodoFetch
	checkErr(err)
	for rows.Next() {
		var id int
		var title string
		var description string
		if err := rows.Scan(&id, &title,&description); err != nil {
			log.Fatal(err)
		}
		todo:= TodoFetch{Id:id,Title:title,Description:description}
		list = append(list,todo)
	}
	return list
}


func checkErr(err error){
	if err!=nil{
		log.Fatal(err)
	}
}