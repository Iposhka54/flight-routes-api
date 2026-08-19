package service

import (
	apperr "flight-routes-api/internal/error"
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

	return airports, nil
}

func (s *AirportService) GetAirport(iataCode string) (model.Airport, error) {
	if err := validateIATACode(iataCode); err != nil {
		return model.Airport{}, err
	}

	airport, err := s.repository.GetAirport(iataCode)
	if err != nil {
		return airport, err
	}

	return airport, nil
}

func (s *AirportService) CreateAirport(airport model.Airport) (model.Airport, error) {
	if err := validateIATACode(airport.IATACode); err != nil {
		return airport, err
	}

	countryLength := len(airport.Country)
	if countryLength < 3 || countryLength > 64 {
		return airport, apperr.ErrInvalidCountry
	}

	nameLength := len(airport.Name)
	if nameLength < 3 || nameLength > 64 {
		return airport, apperr.ErrInvalidName
	}

	id, err := s.repository.CreateAirport(airport)
	if err != nil {
		return airport, err
	}

	airport.ID = id
	return airport, nil
}

func validateIATACode(iataCode string) error {
	if iataCode == "" {
		return apperr.ErrMissingIATACode
	}
	if len(iataCode) != 3 {
		return apperr.ErrInvalidIATACode.WithMessage("Неверный iataCode: %s", iataCode)
	}
	return nil
}
