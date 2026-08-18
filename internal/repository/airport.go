package repository

import (
	"database/sql"
	"flight-routes-api/internal/model"
)

var (
	getAirports      = `SELECT id, iata_code, name, country FROM airport;`
	getAirportByIata = `SELECT id, iata_code, name, country FROM airport
                    	WHERE iata_code = $1;`
	createAirport = `INSERT INTO airport(iata_code, name, country)
					 VALUES ($1, $2, $3) RETURNING id;`
)

type AirportRepository struct {
	db *sql.DB
}

func NewAirportRepository(db *sql.DB) *AirportRepository {
	return &AirportRepository{db: db}
}

func (r *AirportRepository) GetAirports() ([]model.Airport, error) {
	rows, err := r.db.Query(getAirports)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var airports = make([]model.Airport, 0)
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

func (r *AirportRepository) GetAirport(iataCode string) (model.Airport, error) {
	var airport model.Airport
	err := r.db.QueryRow(getAirportByIata, iataCode).Scan(&airport.ID, &airport.IATACode, &airport.Name, &airport.Country)

	if err != nil {
		return model.Airport{}, err
	}

	return airport, nil
}

func (r *AirportRepository) CreateAirport(airport model.Airport) (int, error) {
	var id int
	err := r.db.QueryRow(createAirport,
		airport.IATACode, airport.Name, airport.Country,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}
