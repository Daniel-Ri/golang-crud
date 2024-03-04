package main

import (
    "database/sql"
    "fmt"
    _ "github.com/lib/pq"
	"net/http"
	"time"
)

// handler function for the web server
func handler(w http.ResponseWriter, r *http.Request) {
    // Simulate some processing time
    for i := 0; i < 5; i++ {
        fmt.Fprintf(w, "Processing request...\n")
        // Introducing a small delay to simulate processing time
        // This could be replaced with actual processing logic
        time.Sleep(100 * time.Millisecond)
    }
}

func main() {
    // Establish connection to the PostgreSQL database
    connStr := "user=postgres dbname=test sslmode=disable"
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        panic(err)
    }
    defer db.Close()

    // Attempt to ping the database to check if the connection is successful
    err = db.Ping()
    if err != nil {
        panic(err)
    }

    fmt.Println("Successfully connected to the database!")

    // Perform database operations here
    // For example, you can execute SQL queries using db.Query or db.Exec
	// Register the handler function to handle all requests to "/"
    http.HandleFunc("/", handler)

    // Start the web server in a goroutine
    go func() {
        fmt.Println("Starting server on port 8080...")
        if err := http.ListenAndServe(":8080", nil); err != nil {
            fmt.Printf("Server error: %s\n", err)
        }
    }()

    // This line will execute immediately after starting the server
    fmt.Println("Server started.")

    // Wait indefinitely to keep the program running
    select {}
}
