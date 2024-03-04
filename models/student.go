package models

type Student struct {
	StudentID string `json:"studentId"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Gender    string `json:"gender"`
	Major     string `json:"major"`
}
