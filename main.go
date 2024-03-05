package main

import (
	"daniel/golang-crud/routes"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/lib/pq"
)

// Initialize a global variable for the database connection
var db *sql.DB

func main() {
	// Establish connection to the PostgreSQL database
	connStr := "user=postgres dbname=coursemanagement sslmode=disable"
	var err error
	db, err = CreateConnectionPool("postgres", connStr, 10, 5)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	router := routes.NewRouter(db)

	// Start the server
	fmt.Println("Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}

// CreateConnectionPool creates a connection pool to the database.
func CreateConnectionPool(driverName, dataSourceName string, maxOpenConns, maxIdleConns int) (*sql.DB, error) {
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}

	// Set maximum open connections and maximum idle connections
	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)

	// Ping the database to ensure the connection is valid
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	// Set the maximum lifetime of a connection
	db.SetConnMaxLifetime(time.Minute * 5)

	return db, nil
}
