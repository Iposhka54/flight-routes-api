package service

import (
	"flight-routes-api/internal/model"
	"flight-routes-api/internal/repository"
)

type AirportService struct {
	repository *repository.AirportRepository
}

func NewAirportService(airportRepository *repository.AirportRepository) *AirportService {
	return &AirportService{repository: airportRepository}
}

func (s *AirportService) GetAirports() ([]model.Airport, error) {
	airports, err := s.repository.GetAirports()

	if err != nil {
		return nil, err
	}

	return airports, err
}

func (s *AirportService) GetAirport(iataCode string) (model.Airport, error) {
	airport, err := s.repository.GetAirport(iataCode)

	if err != nil {
		return airport, err
	}

	return airport, nil
}
