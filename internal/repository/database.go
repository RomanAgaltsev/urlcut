package repository

import (
    "context"
    "database/sql"
    "log/slog"
    "time"

    "github.com/RomanAgaltsev/urlcut/internal/interfaces"
    "github.com/RomanAgaltsev/urlcut/internal/model"

    "github.com/pressly/goose/v3"
)

var _ interfaces.Repository = (*DBRepository)(nil)

type DBRepository struct {
    db *sql.DB
}

func NewDBRepository(databaseDSN string) (*DBRepository, error) {
    dbRepository := &DBRepository{}

    if err := dbRepository.migrate(databaseDSN); err != nil {
        return nil, err
    }

    db, err := sql.Open("pgx", databaseDSN)
    if err != nil {
        slog.Error(
            "failed to open DB connection",
            slog.String("error", err.Error()),
        )
    }
    db.SetMaxIdleConns(5)
    db.SetMaxOpenConns(5)
    db.SetConnMaxIdleTime(1 * time.Second)
    db.SetConnMaxLifetime(30 * time.Second)

    dbRepository.db = db

    return dbRepository, nil
}

func (r *DBRepository) migrate(databaseDSN string) error {
    db, err := goose.OpenDBWithDriver("pgx", databaseDSN)
    if err != nil {
        slog.Error(
            "goose: failed to open DB connection",
            slog.String("error", err.Error()),
        )
        return err
    }
    defer func() {
        if err := db.Close(); err != nil {
            slog.Error(
                "goose: failed to close DB connection",
                slog.String("error", err.Error()),
            )
        }
    }()

    if err = goose.SetDialect("postgres"); err != nil {
        slog.Error(
            "goose: failed to set dialect",
            slog.String("error", err.Error()),
        )
        return err
    }

    if err = goose.Up(db, "migrations"); err != nil {
        slog.Error(
            "goose: failed to run migrations",
            slog.String("error", err.Error()),
        )
    }

    return nil
}

func (r *DBRepository) Store(url *model.URL) error {
    return nil
}

func (r *DBRepository) Get(id string) (*model.URL, error) {
    return &model.URL{}, nil
}

func (r *DBRepository) Close() error {
    return nil
}

func (r *DBRepository) Check() error {
    return r.db.PingContext(context.Background())
}
