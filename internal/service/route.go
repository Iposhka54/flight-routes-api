package service

import "flight-routes-api/internal/model"

type Route struct {
	TotalPrice int64
	Flights    []model.Flight
}
