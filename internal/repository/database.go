package repository

import (
    "context"
    "embed"

    //"database/sql"
    "log/slog"
    "time"

    "github.com/RomanAgaltsev/urlcut/internal/interfaces"
    "github.com/RomanAgaltsev/urlcut/internal/model"
    "github.com/RomanAgaltsev/urlcut/internal/repository/queries"

    "github.com/jackc/pgx/v5"
    _ "github.com/jackc/pgx/v5/stdlib"
    "github.com/pressly/goose/v3"
)

var _ interfaces.Repository = (*DBRepository)(nil)

//go:embed migrations/*.sql
var embedMigrations embed.FS

type DBRepository struct {
    conn *pgx.Conn
}

func NewDBRepository(databaseDSN string) (*DBRepository, error) {
    dbRepository := &DBRepository{}

    if err := dbRepository.migrate(databaseDSN); err != nil {
        return nil, err
    }

    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    conn, err := pgx.Connect(ctx, databaseDSN)
    if err != nil {
        slog.Error(
            "failed to open DB connection",
            slog.String("error", err.Error()),
        )
        return nil, err
    }

//    db, err := sql.Open("pgx", databaseDSN)
//    db.SetMaxIdleConns(5)
//    db.SetMaxOpenConns(5)
//    db.SetConnMaxIdleTime(1 * time.Second)
//    db.SetConnMaxLifetime(30 * time.Second)

    dbRepository.conn = conn

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

    ctx, cancel := context.WithTimeout(context.Background(),10 * time.Second)
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
    q := queries.New(r.conn)

    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    _, err := q.CreateURL(ctx, queries.CreateURLParams{
        LongUrl: url.Long,
        BaseUrl: url.Base,
        UrlID: url.ID,
    })
    if err != nil {
        return err
    }

    return nil
}

func (r *DBRepository) Get(id string) (*model.URL, error) {
    q := queries.New(r.conn)

    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    url, err := q.GetURL(ctx, id)
    if err != nil {
        return nil, err
    }

    return &model.URL{
        Long: url.LongUrl,
        Base: url.BaseUrl,
        ID: url.UrlID,
    }, nil
}

func (r *DBRepository) Close() error {
    return r.conn.Close(context.Background())
}

func (r *DBRepository) Check() error {
    return r.conn.Ping(context.Background())
}
