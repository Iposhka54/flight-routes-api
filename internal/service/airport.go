package service

import (
	"context"
	"flight-routes-api/internal/model"
	"flight-routes-api/internal/repository"
)

type AirportService struct {
	repository *repository.AirportRepository
}

func NewAirportService(airportRepository *repository.AirportRepository) *AirportService {
	return &AirportService{repository: airportRepository}
}

func (s *AirportService) GetAirports(ctx context.Context) ([]model.Airport, error) {
	return s.repository.GetAirports(ctx)
}

func (s *AirportService) GetAirport(ctx context.Context, iataCode string) (model.Airport, error) {
	iataCode, err := validateIATACode(iataCode)
	if err != nil {
		return model.Airport{}, err
	}
	return s.repository.GetAirport(ctx, iataCode)
}

func (s *AirportService) CreateAirport(ctx context.Context, airport model.Airport) (model.Airport, error) {
	airport, err := validateAirport(airport)
	if err != nil {
		return airport, err
	}

	id, err := s.repository.CreateAirport(ctx, airport)
	if err != nil {
		return airport, err
	}

	airport.ID = id
	return airport, nil
}
