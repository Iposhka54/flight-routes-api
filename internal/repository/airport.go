package repository

import (
	"database/sql"
	"flight-routes-api/internal/model"
)

type AirportRepository struct {
	db *sql.DB
}

func NewAirportRepository(db *sql.DB) AirportRepository {
	return AirportRepository{db: db}
}

func (h *AirportRepository) GetAirports() ([]model.Airport, error) {
	getAirports := `SELECT * FROM airport`

	rows, err := h.db.Query(getAirports)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var airports []model.Airport

	for rows.Next() {
		var airport model.Airport

		err = rows.Scan(&airport.ID, &airport.IATACode, &airport.Name, &airport.Country)
		if err != nil {
			return nil, err
		}

		airports = append(airports, airport)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return airports, nil
}
