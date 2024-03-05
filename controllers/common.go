package controllers

// Error message struct
type ErrorMessage struct {
	Message string `json:"message"`
}

// Thrown error struct
type ThrownError struct {
	Code    int
	Message string
}
