package model

type Flight struct {
	ID                 int
	OriginAirport      Airport
	DestinationAirport Airport
	Price              int64
}
