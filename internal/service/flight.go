package service

import (
	apperr "flight-routes-api/internal/error"
	"flight-routes-api/internal/model"
	"flight-routes-api/internal/repository"
)

type FlightService struct {
	repository *repository.FlightRepository
}

func NewFlightService(flightRepository *repository.FlightRepository) *FlightService {
	return &FlightService{repository: flightRepository}
}

func (s *FlightService) GetFlights() ([]model.Flight, error) {
	flights, err := s.repository.GetFlights()
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

	flight, err := s.repository.GetFlight(from, to)
	if err != nil {
		return model.Flight{}, err
	}

	return flight, nil
}
