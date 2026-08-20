package service

import (
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
