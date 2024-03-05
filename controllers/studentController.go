package controllers

import (
	"daniel/golang-crud/models"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/lib/pq"
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
		errorMessage := ErrorMessage{Message: "Student Not Found"}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(errorMessage)
		return
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
		log.Fatal(err)
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

	_, err = db.Exec("DELETE FROM students WHERE student_id=$1", params["studentId"])
	if err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(w).Encode(map[string]string{"message": "Student deleted successfully"})
}

// Enroll course
func EnrollCourses(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	w.Header().Set("Content-Type", "application/json")

	defer func() {
		if r := recover(); r != nil {
			if thrownErr, ok := r.(ThrownError); ok {
				w.WriteHeader(thrownErr.Code)
				json.NewEncoder(w).Encode(ErrorMessage{Message: thrownErr.Message})
			} else {
				message := fmt.Sprintf("%v", r)
				json.NewEncoder(w).Encode(ErrorMessage{Message: message})
			}
		}
	}()

	tx, err := db.Begin()
	if err != nil {
		panic(ThrownError{Code: http.StatusInternalServerError, Message: err.Error()})
	}
	defer tx.Rollback() // Rollback transaction if not committed

	params := mux.Vars(r)
	var inputs struct {
		CourseIDs []int `json:"courseIds"`
	}
	_ = json.NewDecoder(r.Body).Decode(&inputs)

	var count int
	err = tx.QueryRow("SELECT COUNT(*) FROM students WHERE student_id = $1", params["studentId"]).Scan(&count)
	if err != nil {
		log.Fatal(err)
	}
	if count == 0 {
		panic(ThrownError{Code: http.StatusNotFound, Message: "Student Not Found"})
	}

	err = tx.QueryRow("SELECT COUNT(*) FROM courses WHERE id = ANY($1)", pq.Array(inputs.CourseIDs)).Scan(&count)
	if err != nil {
		log.Fatal(err)
	}
	if count != len(inputs.CourseIDs) {
		panic(ThrownError{Code: http.StatusNotFound, Message: "Some course IDs Are Not Found"})
	}

	err = tx.QueryRow("SELECT COUNT(*) FROM enrollments WHERE student_id = $1 AND course_id = ANY($2)", params["studentId"], pq.Array(inputs.CourseIDs)).Scan(&count)
	if err != nil {
		log.Fatal(err)
	}
	if count > 0 {
		panic(ThrownError{Code: http.StatusConflict, Message: "You can't enrolled some courses again"})
	}

	err = tx.QueryRow(`
			SELECT COUNT(*)
			FROM (
				SELECT course_id, COUNT(*) AS count_enroll
				FROM enrollments
				WHERE course_id = ANY($1)
				GROUP BY course_id
			) AS subquery 
			JOIN courses ON subquery.course_id = courses.id
			WHERE count_enroll >= max_capacity;
		`, pq.Array(inputs.CourseIDs)).
		Scan(&count)
	if err != nil {
		log.Fatal(err)
	}
	if count > 0 {
		panic(
			ThrownError{Code: http.StatusConflict, Message: "Some courses have reached max capacity"},
		)
	}

	// Insert enrollments to database
	currentDate := time.Now().Format("2006-01-02")
	stmt, err := tx.Prepare("INSERT INTO enrollments (student_id, course_id, enrollment_date) VALUES ($1, $2, $3)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	for _, courseID := range inputs.CourseIDs {
		_, err := stmt.Exec(params["studentId"], courseID, currentDate)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		panic(ThrownError{Code: http.StatusInternalServerError, Message: err.Error()})
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Success enroll courses"})
}
