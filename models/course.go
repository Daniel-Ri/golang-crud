package models

type Course struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	MaxCapacity int    `json:"maxCapacity"`
	Credits     int    `json:"credits"`
}
