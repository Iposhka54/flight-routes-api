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
	getFlightsByOriginIATA = flightSelect + `
		WHERE o.iata_code = $1`
	getFlightsByDestinationIATA = flightSelect + `
		WHERE d.iata_code = $1`
	createFlight = `
		INSERT INTO flight(origin_airport_id, destination_airport_id, price)
		VALUES ($1, $2, $3)
		RETURNING id;`
	updateFlightPrice = `
		UPDATE flight
		SET price = $3
		WHERE origin_airport_id = (SELECT id FROM airport WHERE iata_code = $1)
		  AND destination_airport_id = (SELECT id FROM airport WHERE iata_code = $2)`
)

type FlightRepository struct {
	db *sql.DB
}

func NewFlightRepository(db *sql.DB) *FlightRepository {
	return &FlightRepository{db: db}
}

func (r *FlightRepository) GetFlights() ([]model.Flight, error) {
	return r.queryFlights(getFlights)
}

func (r *FlightRepository) GetFlightsByOriginIATA(originIATA string) ([]model.Flight, error) {
	return r.queryFlights(getFlightsByOriginIATA, originIATA)
}

func (r *FlightRepository) GetFlightsByDestinationIATA(destinationIATA string) ([]model.Flight, error) {
	return r.queryFlights(getFlightsByDestinationIATA, destinationIATA)
}

func (r *FlightRepository) queryFlights(query string, args ...any) ([]model.Flight, error) {
	rows, err := r.db.Query(query, args...)
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

func (r *FlightRepository) UpdateFlightPrice(from, to string, price int64) (model.Flight, error) {
	result, err := r.db.Exec(updateFlightPrice, from, to, price)
	if err != nil {
		return model.Flight{}, wrapDBError(err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return model.Flight{}, wrapDBError(err)
	}
	if affected == 0 {
		return model.Flight{}, apperr.ErrFlightNotFound
	}

	return r.GetFlight(from, to)
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
