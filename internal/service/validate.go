package service

import (
	"strings"

	apperr "flight-routes-api/internal/error"
	"flight-routes-api/internal/model"
)

func normalizeIATACode(code string) string {
	return strings.ToUpper(strings.TrimSpace(code))
}

func validateIATACode(code string) (string, error) {
	code = normalizeIATACode(code)
	if code == "" {
		return "", apperr.ErrMissingIATACode
	}
	if len(code) != 3 || !isAllUpperLatin(code) {
		return "", apperr.ErrInvalidIATACode.WithMessage("Неверный iataCode: %s", code)
	}
	return code, nil
}

func validateRoute(from, to string) (string, string, error) {
	from = normalizeIATACode(from)
	to = normalizeIATACode(to)
	switch {
	case from == "":
		return "", "", apperr.ErrMissingFrom
	case to == "":
		return "", "", apperr.ErrMissingTo
	}

	from, err := validateIATACode(from)
	if err != nil {
		return "", "", err
	}
	to, err = validateIATACode(to)
	if err != nil {
		return "", "", err
	}
	return from, to, nil
}

func validatePrice(price int64) error {
	if price <= 0 {
		return apperr.ErrInvalidPrice
	}
	return nil
}

func validateAirport(airport model.Airport) (model.Airport, error) {
	code, err := validateIATACode(airport.IATACode)
	if err != nil {
		return airport, err
	}
	airport.IATACode = code

	if n := len(airport.Country); n < 3 || n > 64 {
		return airport, apperr.ErrInvalidCountry
	}
	if n := len(airport.Name); n < 3 || n > 64 {
		return airport, apperr.ErrInvalidName
	}
	return airport, nil
}

func validateCreateFlight(origin, destination string, price int64) (string, string, error) {
	origin, err := validateIATACode(origin)
	if err != nil {
		return "", "", err
	}
	destination, err = validateIATACode(destination)
	if err != nil {
		return "", "", err
	}
	if origin == destination {
		return "", "", apperr.ErrSameAirports
	}
	if err := validatePrice(price); err != nil {
		return "", "", err
	}
	return origin, destination, nil
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
