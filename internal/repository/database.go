package repository

import (
    "context"
    "database/sql"
    "log/slog"
    "time"

    "github.com/RomanAgaltsev/urlcut/internal/interfaces"
    "github.com/RomanAgaltsev/urlcut/internal/model"
)

var _ interfaces.Repository = (*DBRepository)(nil)

type DBRepository struct {
    db *sql.DB
}

func NewDBRepository(databaseDSN string) *DBRepository {
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

    return &DBRepository{
        db: db,
    }
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
