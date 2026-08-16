package model

type Airport struct {
	ID       int    `json:"id"`
	IATACode string `json:"iataCode"`
	Name     string `json:"name"`
	Country  string `json:"country"`
}
