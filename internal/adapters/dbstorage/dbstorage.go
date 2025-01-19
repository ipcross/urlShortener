package dbstorage

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	l "github.com/ipcross/urlShortener/internal/logger"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
)

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
	defer dbclose(db)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err = db.PingContext(ctx); err != nil {
		logger.Info("CheckPing", zap.String("error", err.Error()))
		return fmt.Errorf("CheckPing: %w", err)
	}
	return nil
}

func dbclose(db *sql.DB) {
	_ = db.Close()
}
