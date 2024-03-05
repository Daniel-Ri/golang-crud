package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Error message struct
type ErrorMessage struct {
	Message string `json:"message"`
}

// Thrown error struct
type ThrownError struct {
	Code    int
	Message string
}

func CatchPanic(w http.ResponseWriter) {
	if r := recover(); r != nil {
		if thrownErr, ok := r.(ThrownError); ok {
			w.WriteHeader(thrownErr.Code)
			json.NewEncoder(w).Encode(ErrorMessage{Message: thrownErr.Message})
		} else {
			message := fmt.Sprintf("%v", r)
			json.NewEncoder(w).Encode(ErrorMessage{Message: message})
		}
	}
}
