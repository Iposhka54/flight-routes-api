package service

import (
	"cmp"
	"slices"

	apperr "flight-routes-api/internal/error"
	"flight-routes-api/internal/model"
	"flight-routes-api/internal/repository"
)

type FlightService struct {
	flights  *repository.FlightRepository
	airports *repository.AirportRepository
}

func NewFlightService(
	flightRepository *repository.FlightRepository,
	airportRepository *repository.AirportRepository,
) *FlightService {
	return &FlightService{
		flights:  flightRepository,
		airports: airportRepository,
	}
}

func (s *FlightService) GetFlights() ([]model.Flight, error) {
	return s.flights.GetFlights()
}

func (s *FlightService) GetFlight(from, to string) (model.Flight, error) {
	if err := validateRoute(from, to); err != nil {
		return model.Flight{}, err
	}
	return s.flights.GetFlight(from, to)
}

func (s *FlightService) CreateFlight(from, to string, price int64) (model.Flight, error) {
	if err := validateCreateFlight(from, to, price); err != nil {
		return model.Flight{}, err
	}

	originAirport, err := s.airports.GetAirport(from)
	if err != nil {
		return model.Flight{}, err
	}
	destinationAirport, err := s.airports.GetAirport(to)
	if err != nil {
		return model.Flight{}, err
	}

	id, err := s.flights.CreateFlight(originAirport.ID, destinationAirport.ID, price)
	if err != nil {
		return model.Flight{}, err
	}

	return model.Flight{
		ID:                 id,
		OriginAirport:      originAirport,
		DestinationAirport: destinationAirport,
		Price:              price,
	}, nil
}

func (s *FlightService) UpdateFlightPrice(from, to string, price int64) (model.Flight, error) {
	if err := validateRoute(from, to); err != nil {
		return model.Flight{}, err
	}
	if err := validatePrice(price); err != nil {
		return model.Flight{}, err
	}
	return s.flights.UpdateFlightPrice(from, to, price)
}

func (s *FlightService) Search(from, to string) ([]Route, error) {
	if err := validateRoute(from, to); err != nil {
		return nil, err
	}

	outgoing, err := s.flights.GetFlightsByOriginIATA(from)
	if err != nil {
		return nil, err
	}
	incoming, err := s.flights.GetFlightsByDestinationIATA(to)
	if err != nil {
		return nil, err
	}

	incomingByOrigin := make(map[int]model.Flight, len(incoming))
	for _, flight := range incoming {
		incomingByOrigin[flight.OriginAirport.ID] = flight
	}

	routes := make([]Route, 0)

	for _, first := range outgoing {
		if first.DestinationAirport.IATACode == to {
			routes = append(routes, Route{
				TotalPrice: first.Price,
				Flights:    []model.Flight{first},
			})
			continue
		}
		second, ok := incomingByOrigin[first.DestinationAirport.ID]
		if !ok {
			continue
		}
		routes = append(routes, Route{
			TotalPrice: first.Price + second.Price,
			Flights:    []model.Flight{first, second},
		})
	}

	if len(routes) == 0 {
		return nil, apperr.ErrRouteNotFound
	}

	slices.SortFunc(routes, func(a, b Route) int {
		return cmp.Compare(a.TotalPrice, b.TotalPrice)
	})

	return routes, nil
}
