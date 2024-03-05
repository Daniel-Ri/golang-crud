package models

import "time"

type Enrollment struct {
	StudentID      string    `json:"studentId"`
	CourseID       int       `json:"courseId"`
	EnrollmentDate time.Time `json:"enrollmentDate"`
	Grade          float64   `json:"grade"`
}
