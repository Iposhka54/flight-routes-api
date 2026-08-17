package model

type Flight struct {
	ID                   int     `json:"id"`
	OriginAirportID      int     `json:"originAirportId"`
	DestinationAirportID int     `json:"destinationAirportId"`
	Price                float64 `json:"price"`
}
