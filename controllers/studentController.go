package controllers

import (
	"daniel/golang-crud/models"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Get all students
func GetAllStudents(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	w.Header().Set("Content-Type", "application/json")

	rows, err := db.Query("SELECT * FROM students")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var students []models.Student
	for rows.Next() {
		var student models.Student
		err := rows.Scan(&student.StudentID, &student.Name, &student.Email, &student.Gender, &student.Major)
		if err != nil {
			log.Fatal(err)
		}
		students = append(students, student)
	}
	json.NewEncoder(w).Encode(students)
}

// Get one student by ID
func GetStudent(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r) // Get params
	var student models.Student
	err := db.QueryRow("SELECT * FROM students WHERE student_id = $1", params["studentId"]).Scan(&student.StudentID, &student.Name, &student.Email, &student.Gender, &student.Major)
	if err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(w).Encode(student)
}

// Create a new student
func CreateStudent(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	w.Header().Set("Content-Type", "application/json")

	var student models.Student
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
func UpdateStudent(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	var student models.Student
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

	var updatedStudent models.Student
	err = db.QueryRow("SELECT * FROM students WHERE student_id = $1", params["studentId"]).Scan(&updatedStudent.StudentID, &updatedStudent.Name, &updatedStudent.Email, &updatedStudent.Gender, &updatedStudent.Major)
	if err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(w).Encode(updatedStudent)
}

// Delete student data by ID
func DeleteStudent(w http.ResponseWriter, r *http.Request, db *sql.DB) {
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
