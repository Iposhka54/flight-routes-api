package handler

import (
	"flight-routes-api/internal/model"
	"flight-routes-api/internal/service"
	"net/http"
)

type flightResponse struct {
	ID                 int             `json:"id"`
	OriginAirport      airportResponse `json:"originAirport"`
	DestinationAirport airportResponse `json:"destinationAirport"`
	Price              float64         `json:"price"`
}

type FlightHandler struct {
	flightService *service.FlightService
}

func NewFlightHandler(flightService *service.FlightService) *FlightHandler {
	return &FlightHandler{flightService: flightService}
}

func (h *FlightHandler) GetFlights(w http.ResponseWriter, _ *http.Request) error {
	flights, err := h.flightService.GetFlights()
	if err != nil {
		return err
	}

	return writeJSON(w, http.StatusOK, newFlightResponses(flights))
}

func newFlightResponse(flight model.Flight) flightResponse {
	return flightResponse{
		ID:                 flight.ID,
		OriginAirport:      newAirportResponse(flight.OriginAirport),
		DestinationAirport: newAirportResponse(flight.DestinationAirport),
		Price:              float64(flight.Price) / 100,
	}
}

func newFlightResponses(flights []model.Flight) []flightResponse {
	result := make([]flightResponse, 0, len(flights))
	for _, flight := range flights {
		result = append(result, newFlightResponse(flight))
	}
	return result
}
