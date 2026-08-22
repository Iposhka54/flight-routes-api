package service

import (
	apperr "flight-routes-api/internal/error"
	"flight-routes-api/internal/model"
)

func validateIATACode(code string) error {
	if code == "" {
		return apperr.ErrMissingIATACode
	}
	if len(code) != 3 || !isAllUpperLatin(code) {
		return apperr.ErrInvalidIATACode.WithMessage("Неверный iataCode: %s", code)
	}
	return nil
}

func validateIATACodes(codes ...string) error {
	for _, code := range codes {
		if err := validateIATACode(code); err != nil {
			return err
		}
	}
	return nil
}

func validateRoute(from, to string) error {
	switch {
	case from == "":
		return apperr.ErrMissingFrom
	case to == "":
		return apperr.ErrMissingTo
	default:
		return validateIATACodes(from, to)
	}
}

func validatePrice(price int64) error {
	if price <= 0 {
		return apperr.ErrInvalidPrice
	}
	return nil
}

func validateAirport(airport model.Airport) error {
	if err := validateIATACode(airport.IATACode); err != nil {
		return err
	}
	if n := len(airport.Country); n < 3 || n > 64 {
		return apperr.ErrInvalidCountry
	}
	if n := len(airport.Name); n < 3 || n > 64 {
		return apperr.ErrInvalidName
	}
	return nil
}

func validateCreateFlight(origin, destination string, price int64) error {
	if err := validateIATACodes(origin, destination); err != nil {
		return err
	}
	if origin == destination {
		return apperr.ErrSameAirports
	}
	return validatePrice(price)
}

func isAllUpperLatin(s string) bool {
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c < 'A' || c > 'Z' {
			return false
		}
	}
	return true
}
