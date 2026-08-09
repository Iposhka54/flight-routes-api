package sqlite

import (
	"context"
	"database/sql"
	"flight-routes-api/internal/config"
	"fmt"
	"time"
)

func New(cfg config.DatabaseConfig) (*sql.DB, error) {
	dsn := buildDSN(cfg)

	db, err := sql.Open("sqlite3", dsn)
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
