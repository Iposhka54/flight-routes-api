package handler

import (
	"encoding/json"
	apperr "flight-routes-api/internal/error"
	"flight-routes-api/internal/model"
	"flight-routes-api/internal/service"
	"fmt"
	"net/http"
)

type AirportHandler struct {
	airportService *service.AirportService
}

func NewAirportHandler(airportService *service.AirportService) *AirportHandler {
	return &AirportHandler{airportService: airportService}
}

func (h *AirportHandler) GetAirports(w http.ResponseWriter, _ *http.Request) error {
	airports, err := h.airportService.GetAirports()
	if err != nil {
		return err
	}

	return writeJSON(w, http.StatusOK, airports)
}

func (h *AirportHandler) GetAirportByIataCode(w http.ResponseWriter, r *http.Request) error {
	airport, err := h.airportService.GetAirport(r.PathValue("iataCode"))
	if err != nil {
		return err
	}

	return writeJSON(w, http.StatusOK, airport)
}

func (h *AirportHandler) CreateAirport(w http.ResponseWriter, r *http.Request) error {
	var airport model.Airport
	if err := json.NewDecoder(r.Body).Decode(&airport); err != nil {
		return fmt.Errorf("%w: %v", apperr.ErrDecodeJSON, err)
	}

	airport, err := h.airportService.CreateAirport(airport)
	if err != nil {
		return err
	}

	return writeJSON(w, http.StatusOK, airport)
}

func writeJSON(w http.ResponseWriter, status int, data any) error {
	body, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("%w: %v", apperr.ErrEncodeJSON, err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if _, err = w.Write(body); err != nil {
		return fmt.Errorf("%w: %v", apperr.ErrEncodeJSON, err)
	}

	return nil
}
