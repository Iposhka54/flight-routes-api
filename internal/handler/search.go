package handler

import (
	"flight-routes-api/internal/service"
	"net/http"
	"strings"
)

type searchResponse struct {
	Origin      string          `json:"origin"`
	Destination string          `json:"destination"`
	Routes      []routeResponse `json:"routes"`
}

type routeResponse struct {
	TotalPrice float64                `json:"totalPrice"`
	Route      []searchFlightResponse `json:"route"`
}

type searchFlightResponse struct {
	FlightID           int             `json:"flightId"`
	OriginAirport      airportResponse `json:"originAirport"`
	DestinationAirport airportResponse `json:"destinationAirport"`
	Price              float64         `json:"price"`
}

func (h *FlightHandler) Search(w http.ResponseWriter, r *http.Request) error {
	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")

	routes, err := h.flightService.Search(r.Context(), from, to)
	if err != nil {
		return err
	}

	return writeJSON(w, http.StatusOK, newSearchResponse(from, to, routes))
}

func newSearchResponse(from, to string, routes []service.Route) searchResponse {
	items := make([]routeResponse, 0, len(routes))
	for _, route := range routes {
		items = append(items, newRouteResponse(route))
	}

	return searchResponse{
		Origin:      strings.ToUpper(strings.TrimSpace(from)),
		Destination: strings.ToUpper(strings.TrimSpace(to)),
		Routes:      items,
	}
}

func newRouteResponse(route service.Route) routeResponse {
	flights := make([]searchFlightResponse, 0, len(route.Flights))
	for _, flight := range route.Flights {
		flights = append(flights, searchFlightResponse{
			FlightID:           flight.ID,
			OriginAirport:      newAirportResponse(flight.OriginAirport),
			DestinationAirport: newAirportResponse(flight.DestinationAirport),
			Price:              float64(flight.Price) / 100,
		})
	}

	return routeResponse{
		TotalPrice: float64(route.TotalPrice) / 100,
		Route:      flights,
	}
}
