package repository

import (
	"database/sql"
	"flight-routes-api/internal/model"
)

var getFlights = `
	SELECT
		f.id,
		f.price,
		o.id, o.iata_code, o.name, o.country,
		d.id, d.iata_code, d.name, d.country
	FROM flight f
	JOIN airport o ON o.id = f.origin_airport_id
	JOIN airport d ON d.id = f.destination_airport_id;
`

type FlightRepository struct {
	db *sql.DB
}

func NewFlightRepository(db *sql.DB) *FlightRepository {
	return &FlightRepository{db: db}
}

func (r *FlightRepository) GetFlights() ([]model.Flight, error) {
	rows, err := r.db.Query(getFlights)
	if err != nil {
		return nil, wrapDBError(err)
	}
	defer rows.Close()

	flights := make([]model.Flight, 0)
	for rows.Next() {
		flight, err := scanFlight(rows)
		if err != nil {
			return nil, wrapDBError(err)
		}
		flights = append(flights, flight)
	}

	if err = rows.Err(); err != nil {
		return nil, wrapDBError(err)
	}

	return flights, nil
}

func scanFlight(rows *sql.Rows) (model.Flight, error) {
	var flight model.Flight
	err := rows.Scan(
		&flight.ID,
		&flight.Price,
		&flight.OriginAirport.ID, &flight.OriginAirport.IATACode, &flight.OriginAirport.Name, &flight.OriginAirport.Country,
		&flight.DestinationAirport.ID, &flight.DestinationAirport.IATACode, &flight.DestinationAirport.Name, &flight.DestinationAirport.Country,
	)
	if err != nil {
		return model.Flight{}, err
	}

	return flight, nil
}
