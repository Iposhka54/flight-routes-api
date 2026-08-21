package service

import (
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
	flights, err := s.flights.GetFlights()
	if err != nil {
		return nil, err
	}

	return flights, nil
}

func (s *FlightService) GetFlight(from, to string) (model.Flight, error) {
	if from == "" {
		return model.Flight{}, apperr.ErrMissingFrom
	}
	if to == "" {
		return model.Flight{}, apperr.ErrMissingTo
	}
	if err := validateIATACode(from); err != nil {
		return model.Flight{}, err
	}
	if err := validateIATACode(to); err != nil {
		return model.Flight{}, err
	}

	flight, err := s.flights.GetFlight(from, to)
	if err != nil {
		return model.Flight{}, err
	}

	return flight, nil
}

func (s *FlightService) CreateFlight(origin, destination string, price int64) (model.Flight, error) {
	if err := validateIATACode(origin); err != nil {
		return model.Flight{}, err
	}
	if err := validateIATACode(destination); err != nil {
		return model.Flight{}, err
	}
	if origin == destination {
		return model.Flight{}, apperr.ErrSameAirports
	}
	if price <= 0 {
		return model.Flight{}, apperr.ErrInvalidPrice
	}

	originAirport, err := s.airports.GetAirport(origin)
	if err != nil {
		return model.Flight{}, err
	}
	destinationAirport, err := s.airports.GetAirport(destination)
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
	if from == "" {
		return model.Flight{}, apperr.ErrMissingFrom
	}
	if to == "" {
		return model.Flight{}, apperr.ErrMissingTo
	}
	if err := validateIATACode(from); err != nil {
		return model.Flight{}, err
	}
	if err := validateIATACode(to); err != nil {
		return model.Flight{}, err
	}
	if price <= 0 {
		return model.Flight{}, apperr.ErrInvalidPrice
	}

	flight, err := s.flights.UpdateFlightPrice(from, to, price)
	if err != nil {
		return model.Flight{}, err
	}

	return flight, nil
}
