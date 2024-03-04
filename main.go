package main

import (
	"daniel/golang-crud/routes"
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

// Initialize a global variable for the database connection
var db *sql.DB

func main() {
	// Establish connection to the PostgreSQL database
	connStr := "user=postgres dbname=coursemanagement sslmode=disable"
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	router := routes.NewRouter(db)

	// Start the server
	fmt.Println("Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
