package service

import (
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

func (s *FlightService) CreateFlight(origin, destination string, price int64) (model.Flight, error) {
	if err := validateCreateFlight(origin, destination, price); err != nil {
		return model.Flight{}, err
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
	if err := validateRoute(from, to); err != nil {
		return model.Flight{}, err
	}
	if err := validatePrice(price); err != nil {
		return model.Flight{}, err
	}
	return s.flights.UpdateFlightPrice(from, to, price)
}
