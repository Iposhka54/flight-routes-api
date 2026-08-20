package handler

import (
	"encoding/json"
	apperr "flight-routes-api/internal/error"
	"flight-routes-api/internal/model"
	"flight-routes-api/internal/service"
	"math"
	"net/http"
)

type flightResponse struct {
	ID                 int             `json:"id"`
	OriginAirport      airportResponse `json:"originAirport"`
	DestinationAirport airportResponse `json:"destinationAirport"`
	Price              float64         `json:"price"`
}

type createFlightRequest struct {
	OriginIataCode      string  `json:"originIataCode"`
	DestinationIataCode string  `json:"destinationIataCode"`
	Price               float64 `json:"price"`
}

type FlightHandler struct {
	flightService *service.FlightService
}

func NewFlightHandler(flightService *service.FlightService) *FlightHandler {
	return &FlightHandler{flightService: flightService}
}

func (h *FlightHandler) GetFlights(w http.ResponseWriter, r *http.Request) error {
	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")
	if from == "" && to == "" {
		flights, err := h.flightService.GetFlights()
		if err != nil {
			return err
		}

		return writeJSON(w, http.StatusOK, newFlightResponses(flights))
	}

	flight, err := h.flightService.GetFlight(from, to)
	if err != nil {
		return err
	}

	return writeJSON(w, http.StatusOK, newFlightResponse(flight))
}

func (h *FlightHandler) CreateFlight(w http.ResponseWriter, r *http.Request) error {
	var req createFlightRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return apperr.ErrDecodeJSON.Wrap(err)
	}

	flight, err := h.flightService.CreateFlight(
		req.OriginIataCode,
		req.DestinationIataCode,
		rublesToKopecks(req.Price),
	)
	if err != nil {
		return err
	}

	return writeJSON(w, http.StatusCreated, newFlightResponse(flight))
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

func rublesToKopecks(rubles float64) int64 {
	return int64(math.Round(rubles * 100))
}
