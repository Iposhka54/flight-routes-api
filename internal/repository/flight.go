package repository

import (
	"database/sql"
	"errors"
	apperr "flight-routes-api/internal/error"
	"flight-routes-api/internal/model"
	"flight-routes-api/internal/sqlite"
)

var flightSelect = `
	SELECT
		f.id,
		f.price,
		o.id, o.iata_code, o.name, o.country,
		d.id, d.iata_code, d.name, d.country
	FROM flight f
	JOIN airport o ON o.id = f.origin_airport_id
	JOIN airport d ON d.id = f.destination_airport_id
`

var (
	getFlights       = flightSelect
	getFlightByRoute = flightSelect + `
		WHERE o.iata_code = $1 AND d.iata_code = $2`
	createFlight = `
		INSERT INTO flight(origin_airport_id, destination_airport_id, price)
		VALUES ($1, $2, $3)
		RETURNING id;`
)

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

func (r *FlightRepository) GetFlight(from, to string) (model.Flight, error) {
	var flight model.Flight
	err := r.db.QueryRow(getFlightByRoute, from, to).Scan(
		&flight.ID,
		&flight.Price,
		&flight.OriginAirport.ID, &flight.OriginAirport.IATACode, &flight.OriginAirport.Name, &flight.OriginAirport.Country,
		&flight.DestinationAirport.ID, &flight.DestinationAirport.IATACode, &flight.DestinationAirport.Name, &flight.DestinationAirport.Country,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Flight{}, apperr.ErrFlightNotFound
		}
		return model.Flight{}, wrapDBError(err)
	}

	return flight, nil
}

func (r *FlightRepository) CreateFlight(originAirportID, destinationAirportID int, price int64) (int, error) {
	var id int
	err := r.db.QueryRow(createFlight, originAirportID, destinationAirportID, price).Scan(&id)
	if err != nil {
		if sqlite.IsUniqueConstraint(err) {
			return 0, apperr.ErrFlightAlreadyExists
		}
		return 0, wrapDBError(err)
	}

	return id, nil
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
