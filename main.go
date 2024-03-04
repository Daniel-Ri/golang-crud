package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type Student struct {
	StudentID string `json:"studentId"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Gender    string `json:"gender"`
	Major     string `json:"major"`
}

// Error message struct
type ErrorMessage struct {
	Message string `json:"message"`
}

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

	// Initialize the router
	router := mux.NewRouter()

	// Define endpoints
	router.HandleFunc("/students", getAllStudents).Methods("GET")
	router.HandleFunc("/students/{studentId}", getStudent).Methods("GET")
	router.HandleFunc("/students", createStudent).Methods("POST")
	router.HandleFunc("/students/{studentId}", updateStudent).Methods("PUT")
	router.HandleFunc("/students/{studentId}", deleteStudent).Methods("DELETE")

	// Start the server
	fmt.Println("Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}

// Get all students
func getAllStudents(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	rows, err := db.Query("SELECT * FROM students")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var students []Student
	for rows.Next() {
		var student Student
		err := rows.Scan(&student.StudentID, &student.Name, &student.Email, &student.Gender, &student.Major)
		if err != nil {
			log.Fatal(err)
		}
		students = append(students, student)
	}
	json.NewEncoder(w).Encode(students)
}

// Get one student by ID
func getStudent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r) // Get params
	var student Student
	err := db.QueryRow("SELECT * FROM students WHERE student_id = $1", params["studentId"]).Scan(&student.StudentID, &student.Name, &student.Email, &student.Gender, &student.Major)
	if err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(w).Encode(student)
}

// Create a new student
func createStudent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var student Student
	_ = json.NewDecoder(r.Body).Decode(&student)

	// Check if student_id already exists
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM students WHERE student_id = $1", student.StudentID).Scan(&count)
	if err != nil {
		log.Fatal(err)
	}
	if count > 0 {
		// Return error response if student_id already exists
		errorMessage := ErrorMessage{Message: "Student ID already exists"}
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(errorMessage)
		return
	}

	_, err = db.Exec("INSERT INTO students (student_id, name, email, gender, major) VALUES ($1, $2, $3, $4, $5)", student.StudentID, student.Name, student.Email, student.Gender, student.Major)
	if err != nil {
		log.Fatal((err))
	}
	json.NewEncoder(w).Encode(student)
}

// Update student data by ID
func updateStudent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	var student Student
	_ = json.NewDecoder(r.Body).Decode(&student)

	// Check if the student exists
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM students WHERE student_id = $1", params["studentId"]).Scan(&count)
	if err != nil {
		log.Fatal(err)
	}
	if count == 0 {
		errorMessage := ErrorMessage{Message: "Student not found"}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(errorMessage)
		return
	}

	// Prepare the UPDATE statement
	updateQuery := "UPDATE students SET "
	var values []interface{}
	paramCount := 1 // Start with $1 as the first parameter
	if student.Name != "" || r.FormValue("name") != "" {
		updateQuery += fmt.Sprintf("name = $%d, ", paramCount)
		values = append(values, student.Name)
		paramCount++
	}
	if student.Email != "" || r.FormValue("email") != "" {
		updateQuery += fmt.Sprintf("email = $%d, ", paramCount)
		values = append(values, student.Email)
		paramCount++
	}
	if student.Gender != "" || r.FormValue("gender") != "" {
		updateQuery += fmt.Sprintf("gender = $%d, ", paramCount)
		values = append(values, student.Gender)
		paramCount++
	}
	if student.Major != "" || r.FormValue("major") != "" {
		updateQuery += fmt.Sprintf("major = $%d, ", paramCount)
		values = append(values, student.Major)
		paramCount++
	}

	// Remove the trailing comma and space
	updateQuery = updateQuery[:len(updateQuery)-2]

	// Add the WHERE clause for the student_id
	updateQuery += " WHERE student_id = $"
	updateQuery += fmt.Sprintf("%d", paramCount)
	values = append(values, params["studentId"])

	// Execute the UPDATE statement
	_, err = db.Exec(updateQuery, values...)
	if err != nil {
		log.Fatal(err)
	}

	var updatedStudent Student
	err = db.QueryRow("SELECT * FROM students WHERE student_id = $1", params["studentId"]).Scan(&updatedStudent.StudentID, &updatedStudent.Name, &updatedStudent.Email, &updatedStudent.Gender, &updatedStudent.Major)
	if err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(w).Encode(updatedStudent)
}

// Delete student data by ID
func deleteStudent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	// Check if the student exists
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM students WHERE student_id = $1", params["studentId"]).Scan(&count)
	if err != nil {
		log.Fatal(err)
	}
	if count == 0 {
		errorMessage := ErrorMessage{Message: "Student Not Found"}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(errorMessage)
		return
	}

	_, err = db.Exec("DELETE FROM students WHERE student_id=$1", params["id"])
	if err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(w).Encode(map[string]string{"message": "Student deleted successfully"})
}
