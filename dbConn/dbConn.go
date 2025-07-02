package dbconn
import (
	"database/sql"
	"log"
	"os"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// This function is used to create a connection of our server with the database
func ConnToDb()*sql.DB{
	var err error
	err = godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}
	connStr:= os.Getenv("CONNECTION_STRING")
	var db *sql.DB
	db,err = sql.Open("postgres",connStr)
	checkErr(err)
	// createTable(db)
	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
	return db
}

func checkErr(err error){
	if err!=nil{
		log.Fatal(err)
	}
}