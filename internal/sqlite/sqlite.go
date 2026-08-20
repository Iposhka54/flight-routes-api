package sqlite

import (
	"context"
	"database/sql"
	"flight-routes-api/internal/config"
	"fmt"
	"time"

	_ "modernc.org/sqlite"
)

func New(cfg config.DatabaseConfig) (*sql.DB, error) {
	dsn := buildDSN(cfg)

	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetMaxOpenConns(cfg.MaxOpenConns)

	if cfg.ConnMaxLifetimeSeconds > 0 {
		db.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetimeSeconds) * time.Second)
	} else {
		db.SetConnMaxLifetime(0)
	}

	if err = pingWithTimeout(db, cfg.QueryTimeoutSeconds); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

func buildDSN(cfg config.DatabaseConfig) string {
	dsn := fmt.Sprintf("%s?mode=%s", cfg.Path, cfg.Mode)

	dsn += fmt.Sprintf("&cache_size=%d", cfg.CacheSize*1024)
	dsn += fmt.Sprintf("&synchronous=%s", cfg.SyncMode)

	dsn += fmt.Sprintf("&busy_timeout=%d", cfg.BusyTimeoutSeconds*1000)

	dsn += "&_journal_mode=WAL"
	dsn += "&_foreign_keys=ON"
	dsn += "&_temp_store=MEMORY"
	dsn += "&_auto_vacuum=INCREMENTAL"

	return dsn
}

func pingWithTimeout(db *sql.DB, timeoutSec int) error {
	if timeoutSec <= 0 {
		timeoutSec = 5
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutSec)*time.Second)
	defer cancel()

	return db.PingContext(ctx)
}

func InitDB(db *sql.DB) error {
	airportTable := `CREATE TABLE IF NOT EXISTS airport (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    iata_code VARCHAR(3) UNIQUE NOT NULL,
    name VARCHAR NOT NULL,
    country VARCHAR NOT NULL
    );`

	if _, err := db.Exec(airportTable); err != nil {
		return fmt.Errorf("failed to create airport table: %w", err)
	}

	flightsTable := `
	CREATE TABLE IF NOT EXISTS flight (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		origin_airport_id INTEGER NOT NULL,
		destination_airport_id INTEGER NOT NULL,
		price INTEGER NOT NULL,
		FOREIGN KEY (origin_airport_id) REFERENCES airport(id),
		FOREIGN KEY (destination_airport_id) REFERENCES airport(id),
		UNIQUE(origin_airport_id, destination_airport_id)
	);`

	if _, err := db.Exec(flightsTable); err != nil {
		return fmt.Errorf("failed to create flights table: %w", err)
	}

	return nil
}
