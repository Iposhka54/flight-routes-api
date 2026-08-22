package handler

import (
	"encoding/json"
	apperr "flight-routes-api/internal/error"
	"flight-routes-api/internal/model"
	"flight-routes-api/internal/service"
	"net/http"
)

type airportResponse struct {
	ID       int    `json:"id"`
	IATACode string `json:"iataCode"`
	Name     string `json:"name"`
	Country  string `json:"country"`
}

type createAirportRequest struct {
	IATACode string `json:"iataCode"`
	Name     string `json:"name"`
	Country  string `json:"country"`
}

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

	return writeJSON(w, http.StatusOK, newAirportResponses(airports))
}

func (h *AirportHandler) GetAirportByIataCode(w http.ResponseWriter, r *http.Request) error {
	airport, err := h.airportService.GetAirport(r.PathValue("iataCode"))
	if err != nil {
		return err
	}

	return writeJSON(w, http.StatusOK, newAirportResponse(airport))
}

func (h *AirportHandler) CreateAirport(w http.ResponseWriter, r *http.Request) error {
	var req createAirportRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return apperr.ErrDecodeJSON.Wrap(err)
	}

	airport, err := h.airportService.CreateAirport(req.toModel())
	if err != nil {
		return err
	}

	return writeJSON(w, http.StatusCreated, newAirportResponse(airport))
}

func newAirportResponse(airport model.Airport) airportResponse {
	return airportResponse{
		ID:       airport.ID,
		IATACode: airport.IATACode,
		Name:     airport.Name,
		Country:  airport.Country,
	}
}

func newAirportResponses(airports []model.Airport) []airportResponse {
	result := make([]airportResponse, 0, len(airports))
	for _, airport := range airports {
		result = append(result, newAirportResponse(airport))
	}
	return result
}

func (r createAirportRequest) toModel() model.Airport {
	return model.Airport{
		IATACode: r.IATACode,
		Name:     r.Name,
		Country:  r.Country,
	}
}
