package routes

import (
	"daniel/golang-crud/controllers"
	"database/sql"
	"net/http"

	"github.com/gorilla/mux"
)

func NewRouter(db *sql.DB) *mux.Router {
	// Initialize the router
	router := mux.NewRouter()

	// Define endpoints
	router.HandleFunc("/students", func(w http.ResponseWriter, r *http.Request) {
		controllers.GetAllStudents(w, r, db)
	}).Methods("GET")
	router.HandleFunc("/students/{studentId}", func(w http.ResponseWriter, r *http.Request) {
		controllers.GetStudent(w, r, db)
	}).Methods("GET")
	router.HandleFunc("/students", func(w http.ResponseWriter, r *http.Request) {
		controllers.CreateStudent(w, r, db)
	}).Methods("POST")
	router.HandleFunc("/students/{studentId}", func(w http.ResponseWriter, r *http.Request) {
		controllers.UpdateStudent(w, r, db)
	}).Methods("PUT")
	router.HandleFunc("/students/{studentId}", func(w http.ResponseWriter, r *http.Request) {
		controllers.DeleteStudent(w, r, db)
	}).Methods("DELETE")

	router.HandleFunc("/courses", func(w http.ResponseWriter, r *http.Request) {
		controllers.GetAllCourses(w, r, db)
	}).Methods("GET")

	return router
}
