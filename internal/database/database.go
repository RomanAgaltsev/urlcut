// Пакет database реализует функционал соединения с БД и миграции.
package database

import (
	"context"
	"database/sql"
	"log/slog"
	"time"

	"github.com/pressly/goose/v3"

	"github.com/RomanAgaltsev/urlcut/migrations"
)

// NewConnection создает новое соединение с базой данных.
// Выполняются следующие действия:
//  - открывается соединение с БД с установкой параметров соединения
//  - выполняется пинг БД
//  - запускаются миграции
func NewConnection(ctx context.Context, driver string, databaseDSN string) (*sql.DB, error) {
	// Открываем соединение
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

	// Делаем пинг
	if err = db.PingContext(ctx); err != nil {
		slog.Error("failed to ping DB", slog.String("error", err.Error()))
		return nil, err
	}

	// Запускаем миграции
	Migrate(ctx, databaseDSN)

	return db, err
}

// Migrate выполняет миграции базы данных.
func Migrate(ctx context.Context, databaseDSN string) {
	// Тут открываем своё соединение
	db, err := goose.OpenDBWithDriver("pgx", databaseDSN)
	if err != nil {
		slog.Error("goose: failed to open DB connection", slog.String("error", err.Error()))
	}
	// Откладываем закрытие соединения
	defer func() {
		if err := db.Close(); err != nil {
			slog.Error("goose: failed to close DB connection", slog.String("error", err.Error()))
		}
	}()

	// Устанавливаем диалект
	if err = goose.SetDialect("postgres"); err != nil {
		slog.Error("goose: failed to set dialect", slog.String("error", err.Error()))
	}

	// Устанавливаем папку с файлами миграции
	goose.SetBaseFS(migrations.Migrations)

	// Накатываем миграции
	if err = goose.UpContext(ctx, db, "."); err != nil {
		slog.Error("goose: failed to run migrations", slog.String("error", err.Error()))
	}
}
