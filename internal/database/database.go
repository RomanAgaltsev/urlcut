package database

import (
	"context"
	"database/sql"
	"log/slog"
	"time"

	"github.com/RomanAgaltsev/urlcut/migrations"

	"github.com/pressly/goose/v3"
)

func NewConnection(ctx context.Context, driver string, databaseDSN string) (*sql.DB, error) {
	db, err := sql.Open(driver, databaseDSN)
	if err != nil {
		slog.Error("failed to open DB connection", slog.String("error", err.Error()))
		return nil, err
	}

	// Устанавливаем параметры соединений
	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(5)
	db.SetConnMaxIdleTime(1 * time.Second)
	db.SetConnMaxLifetime(30 * time.Second)

	if err = db.PingContext(ctx); err != nil {
		slog.Error("failed to ping DB", slog.String("error", err.Error()))
		return nil, err
	}

	Migrate(ctx, databaseDSN)

	return db, err
}

func Migrate(ctx context.Context, databaseDSN string) {
	db, err := goose.OpenDBWithDriver("pgx", databaseDSN)
	if err != nil {
		slog.Error("goose: failed to open DB connection", slog.String("error", err.Error()))
	}
	defer func() {
		if err := db.Close(); err != nil {
			slog.Error("goose: failed to close DB connection", slog.String("error", err.Error()))
		}
	}()

	if err = goose.SetDialect("postgres"); err != nil {
		slog.Error("goose: failed to set dialect", slog.String("error", err.Error()))
	}

	goose.SetBaseFS(migrations.Migrations)

	if err = goose.UpContext(ctx, db, "."); err != nil {
		slog.Error("goose: failed to run migrations", slog.String("error", err.Error()))
	}
}
