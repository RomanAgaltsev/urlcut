package repository

import (
    "context"
    "database/sql"
    "embed"
    "log/slog"
    "time"

    "github.com/RomanAgaltsev/urlcut/internal/interfaces"
    "github.com/RomanAgaltsev/urlcut/internal/model"
    "github.com/RomanAgaltsev/urlcut/internal/repository/queries"

    _ "github.com/jackc/pgx/v5/stdlib"
    "github.com/pressly/goose/v3"
)

var _ interfaces.Repository = (*DBRepository)(nil)

//go:embed migrations/*.sql
var embedMigrations embed.FS

type DBRepository struct {
    db *sql.DB
    *queries.Queries
}

func NewDBRepository(databaseDSN string) (*DBRepository, error) {
    db, err := sql.Open("pgx", databaseDSN)
    if err != nil {
        slog.Error(
            "failed to open DB connection",
            slog.String("error", err.Error()),
        )
        return nil, err
    }

    db.SetMaxIdleConns(5)
    db.SetMaxOpenConns(5)
    db.SetConnMaxIdleTime(1 * time.Second)
    db.SetConnMaxLifetime(30 * time.Second)

    dbRepository := &DBRepository{
        db:      db,
        Queries: queries.New(db),
    }

    if err := dbRepository.migrate(databaseDSN); err != nil {
        return nil, err
    }

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

    goose.SetBaseFS(embedMigrations)

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    if err = goose.UpContext(ctx, db, "migrations"); err != nil {
        slog.Error(
            "goose: failed to run migrations",
            slog.String("error", err.Error()),
        )
    }

    return nil
}

func (r *DBRepository) Store(url *model.URL) error {
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    _, err := r.CreateURL(ctx, queries.CreateURLParams{
        LongUrl: url.Long,
        BaseUrl: url.Base,
        UrlID:   url.ID,
    })
    if err != nil {
        return err
    }

    return nil
}

func (r *DBRepository) Get(id string) (*model.URL, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    url, err := r.GetURL(ctx, id)
    if err != nil {
        return nil, err
    }

    return &model.URL{
        Long: url.LongUrl,
        Base: url.BaseUrl,
        ID:   url.UrlID,
    }, nil
}

func (r *DBRepository) Close() error {
    return r.db.Close()
}

func (r *DBRepository) Check() error {
    return r.db.Ping()
}
