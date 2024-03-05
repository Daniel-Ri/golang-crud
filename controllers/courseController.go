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

// Get all courses
func GetAllCourses(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	w.Header().Set("Content-Type", "application/json")

	rows, err := db.Query("SELECT * FROM courses")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var courses []models.Course
	for rows.Next() {
		var course models.Course
		err := rows.Scan(&course.ID, &course.Name, &course.MaxCapacity, &course.Credits)
		if err != nil {
			log.Fatal(err)
		}
		courses = append(courses, course)
	}
	json.NewEncoder(w).Encode(courses)
}

// Get one course by ID
func GetCourse(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r) // Get params
	var course models.Course
	err := db.QueryRow("SELECT * FROM courses WHERE id = $1", params["courseId"]).Scan(&course.ID, &course.Name, &course.MaxCapacity, &course.Credits)
	if err != nil {
		errorMessage := ErrorMessage{Message: "Course Not Found"}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(errorMessage)
		return
	}

	json.NewEncoder(w).Encode(course)
}

// Create a new course
func CreateCourse(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	w.Header().Set("Content-Type", "application/json")

	var course models.Course
	_ = json.NewDecoder(r.Body).Decode(&course)

	// Check if course name already exists
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM courses WHERE name = $1", course.Name).Scan(&count)
	if err != nil {
		log.Fatal(err)
	}
	if count > 0 {
		// Return error response if course name already exists
		errorMessage := ErrorMessage{Message: "Course name already exists"}
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(errorMessage)
		return
	}

	_, err = db.Exec("INSERT INTO courses (name, max_capacity, credits) VALUES ($1, $2, $3)", course.Name, course.MaxCapacity, course.Credits)
	if err != nil {
		log.Fatal(err)
	}

	var createdCourse models.Course
	err = db.
		QueryRow("SELECT * FROM courses WHERE name = $1", course.Name).
		Scan(&createdCourse.ID, &createdCourse.Name, &createdCourse.MaxCapacity, &createdCourse.Credits)
	if err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(w).Encode(createdCourse)
}

// Update course data by ID
func UpdateCourse(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	var course models.Course
	_ = json.NewDecoder(r.Body).Decode(&course)

	// Check if the course exists
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM courses where id = $1", params["courseId"]).Scan(&count)
	if err != nil {
		log.Fatal(err)
	}
	if count == 0 {
		errorMessage := ErrorMessage{Message: "Course Not Found"}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(errorMessage)
		return
	}

	// Prepare the UPDATE statement
	updateQuery := "UPDATE courses SET "
	var values []interface{}
	paramCount := 1 // Start with $1 as the first parameter
	if course.Name != "" || r.FormValue("name") != "" {
		// Check if course name already exists
		var count int
		err := db.QueryRow("SELECT COUNT(*) FROM courses WHERE name = $1 AND id <> $2", course.Name, params["courseId"]).Scan(&count)
		if err != nil {
			log.Fatal(err)
		}
		if count > 0 {
			// Return error response if course name already exists
			errorMessage := ErrorMessage{Message: "Course name already exists"}
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(errorMessage)
			return
		}

		updateQuery += fmt.Sprintf("name = $%d, ", paramCount)
		values = append(values, course.Name)
		paramCount++
	}
	if course.MaxCapacity != 0 || r.FormValue("maxCapacity") != "" {
		updateQuery += fmt.Sprintf("max_capacity = $%d, ", paramCount)
		values = append(values, course.MaxCapacity)
		paramCount++
	}
	if course.Credits != 0 || r.FormValue("credits") != "" {
		updateQuery += fmt.Sprintf("credits = $%d, ", paramCount)
		values = append(values, course.Credits)
		paramCount++
	}

	// Remove the trailing comma and space
	updateQuery = updateQuery[:len(updateQuery)-2]

	// Add the WHERE clause for the id
	updateQuery += fmt.Sprintf(" WHERE id = $%d", paramCount)
	values = append(values, params["courseId"])

	// Execute the UPDATE statement
	_, err = db.Exec(updateQuery, values...)
	if err != nil {
		log.Fatal(err)
	}

	var updatedCourse models.Course
	err = db.
		QueryRow("SELECT * FROM courses WHERE id = $1", params["courseId"]).
		Scan(&updatedCourse.ID, &updatedCourse.Name, &updatedCourse.MaxCapacity, &updatedCourse.Credits)
	if err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(w).Encode(updatedCourse)
}

// Delete course data by ID
func DeleteCourse(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	// Check if the course exists
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM courses WHERE id = $1", params["courseId"]).Scan(&count)
	if err != nil {
		log.Fatal(err)
	}
	if count == 0 {
		errorMessage := ErrorMessage{Message: "Course Not Found"}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(errorMessage)
		return
	}

	_, err = db.Exec("DELETE FROM courses WHERE id=$1", params["courseId"])
	if err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(w).Encode(map[string]string{"message": "Course deleted successfully"})
}
