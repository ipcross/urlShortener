package dbstorage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	l "github.com/ipcross/urlShortener/internal/logger"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
)

const (
	timeout = 20
)

type Event struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

func InsertRecord(dsn string, event *Event) error {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return fmt.Errorf("failed to open a connection to the DB: %w", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("failed to properly close the DB connection: %v", err)
		}
	}()
	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
	defer cancel()
	_, err = db.ExecContext(ctx,
		"insert into url_table(short_url, original_url) values ($1, $2)", event.ShortURL, event.OriginalURL)
	if err != nil {
		return fmt.Errorf("failed to insert data: %w", err)
	}
	return nil
}

func LoadDBData(dsn string) ([]Event, error) {
	if len(dsn) == 0 {
		return nil, errors.New("passed DSN is empty")
	}
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open a connection to the DB: %w", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("failed to properly close the DB connection: %v", err)
		}
	}()
	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
	defer cancel()
	if err := createSchema(ctx, db); err != nil {
		return nil, fmt.Errorf("failed to create the DB schema: %w", err)
	}
	rows, err := db.QueryContext(ctx, "SELECT short_url, original_url FROM url_table")
	if err != nil {
		return nil, fmt.Errorf("failed to load data from DB: %w", err)
	}
	defer func() {
		_ = rows.Close()
		_ = rows.Err()
	}()
	events := make([]Event, 0)
	for rows.Next() {
		event := Event{}
		err = rows.Scan(&event.ShortURL, &event.OriginalURL)
		if err != nil {
			return nil, fmt.Errorf("failed to scan: %w", err)
		}
		events = append(events, event)
	}
	return events, nil
}

func CheckPing(dsn string) error {
	logger, err := l.Initialize("Info")
	if err != nil {
		return fmt.Errorf("logger Initialize: %w", err)
	}

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		logger.Info("CheckPing", zap.String("error", err.Error()))
		return fmt.Errorf("sql.Open: %w", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("failed to properly close the DB connection: %v", err)
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err = db.PingContext(ctx); err != nil {
		logger.Info("CheckPing", zap.String("error", err.Error()))
		return fmt.Errorf("CheckPing: %w", err)
	}
	return nil
}

func createSchema(ctx context.Context, db *sql.DB) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to start a transaction: %w", err)
	}
	defer func() {
		if err := tx.Rollback(); err != nil {
			if !errors.Is(err, sql.ErrTxDone) {
				log.Printf("failed to rollback the transaction: %v", err)
			}
		}
	}()

	stmt := "CREATE TABLE IF NOT EXISTS url_table(short_url TEXT primary key, original_url TEXT)"
	if _, err := tx.ExecContext(ctx, stmt); err != nil {
		return fmt.Errorf("failed to execute statement `%s`: %w", stmt, err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit the transaction: %w", err)
	}
	return nil
}
