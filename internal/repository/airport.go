package repository

import (
	"context"
	"database/sql"
	"errors"
	apperr "flight-routes-api/internal/error"
	"flight-routes-api/internal/model"
	"flight-routes-api/internal/sqlite"
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

func (r *AirportRepository) GetAirports(ctx context.Context) ([]model.Airport, error) {
	rows, err := r.db.QueryContext(ctx, getAirports)
	if err != nil {
		return nil, wrapDBError(err)
	}
	defer rows.Close()

	var airports = make([]model.Airport, 0)
	for rows.Next() {
		var airport model.Airport

		err = rows.Scan(&airport.ID, &airport.IATACode, &airport.Name, &airport.Country)
		if err != nil {
			return nil, wrapDBError(err)
		}

		airports = append(airports, airport)
	}

	if err = rows.Err(); err != nil {
		return nil, wrapDBError(err)
	}

	return airports, nil
}

func (r *AirportRepository) GetAirport(ctx context.Context, iataCode string) (model.Airport, error) {
	var airport model.Airport
	err := r.db.QueryRowContext(ctx, getAirportByIata, iataCode).Scan(&airport.ID, &airport.IATACode, &airport.Name, &airport.Country)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Airport{}, apperr.ErrAirportNotFound
		}
		return model.Airport{}, wrapDBError(err)
	}

	return airport, nil
}

func (r *AirportRepository) CreateAirport(ctx context.Context, airport model.Airport) (int, error) {
	var id int
	err := r.db.QueryRowContext(ctx, createAirport,
		airport.IATACode, airport.Name, airport.Country,
	).Scan(&id)
	if err != nil {
		if sqlite.IsUniqueConstraint(err) {
			return 0, apperr.ErrAirportAlreadyExists
		}
		return 0, wrapDBError(err)
	}

	return id, nil
}

func wrapDBError(err error) error {
	return apperr.ErrInternal.Wrap(err)
}
