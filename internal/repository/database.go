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
    q  *queries.Queries
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

    var q *queries.Queries

    q, err = queries.Prepare(context.Background(), db)
    if err != nil {
        q = queries.New(db)
    }

    dbRepository := &DBRepository{
        db: db,
        q:  q,
    }

    if err := dbRepository.bootstrap(databaseDSN); err != nil {
        return nil, err
    }

    return dbRepository, nil
}

func (r *DBRepository) bootstrap(databaseDSN string) error {
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

func (r *DBRepository) Store(urls []*model.URL) error {
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    tx, err := r.db.Begin()
    if err != nil {
        return err
    }
    defer func() { _ = tx.Rollback() }()

    qtx := r.q.WithTx(tx)

    for _, url := range urls {
        _, err := qtx.StoreURL(ctx, queries.StoreURLParams{
            LongUrl: url.Long,
            BaseUrl: url.Base,
            UrlID:   url.ID,
        })
        if err != nil {
            return err
        }
    }

    return tx.Commit()
}

func (r *DBRepository) Get(id string) (*model.URL, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    url, err := r.q.GetURL(ctx, id)
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
