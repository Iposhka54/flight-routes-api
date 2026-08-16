package handler

import (
	"encoding/json"
	"flight-routes-api/internal/service"
	"net/http"
)

type AirportHandler struct {
	airportService service.AirportService
}

func NewAirportHandler(airportService service.AirportService) AirportHandler {
	return AirportHandler{airportService: airportService}

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
