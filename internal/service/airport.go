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
	return s.repository.GetAirports()
}

func (s *AirportService) GetAirport(iataCode string) (model.Airport, error) {
	if err := validateIATACode(iataCode); err != nil {
		return model.Airport{}, err
	}
	return s.repository.GetAirport(iataCode)
}

func (s *AirportService) CreateAirport(airport model.Airport) (model.Airport, error) {
	if err := validateAirport(airport); err != nil {
		return airport, err
	}

	id, err := s.repository.CreateAirport(airport)
	if err != nil {
		return airport, err
	}

	airport.ID = id
	return airport, nil
}
