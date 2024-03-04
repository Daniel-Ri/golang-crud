package controllers

import (
	"daniel/golang-crud/models"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

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
		err := rows.Scan(&course.ID, &course.Name, &course.Credits)
		if err != nil {
			log.Fatal(err)
		}
		courses = append(courses, course)
	}
	json.NewEncoder(w).Encode(courses)
}
