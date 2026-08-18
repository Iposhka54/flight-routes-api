package handler

import (
	"encoding/json"
	"flight-routes-api/internal/service"
	"net/http"
)

type AirportHandler struct {
	airportService *service.AirportService
}

func NewAirportHandler(airportService *service.AirportService) *AirportHandler {
	return &AirportHandler{airportService: airportService}
}

func (h *AirportHandler) GetAirports(w http.ResponseWriter, _ *http.Request) {
	airports, err := h.airportService.GetAirports()
	if err != nil {
		http.Error(w, "Failed to get airports", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(airports); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *AirportHandler) GetAirportByIataCode(w http.ResponseWriter, r *http.Request) {
	iataCode := r.PathValue("iataCode")

	if iataCode == "" {
		h.respondError(w, http.StatusBadRequest, "missing airport iataCode")
		return
	}

	if len(iataCode) != 3 {
		h.respondError(w, http.StatusBadRequest, "missing iataCode's format")
		return
	}

	airport, err := h.airportService.GetAirport(iataCode)
	if err != nil {
		http.Error(w, "Failed to get airport", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(airport); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *AirportHandler) respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *AirportHandler) respondError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{
		"error": message,
	})
}
