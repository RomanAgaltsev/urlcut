package repository

import (
    "context"
    "database/sql"
    "embed"
    "errors"
    "log/slog"
    "time"

    "github.com/RomanAgaltsev/urlcut/internal/interfaces"
    "github.com/RomanAgaltsev/urlcut/internal/model"
    "github.com/RomanAgaltsev/urlcut/internal/repository/queries"

    "github.com/cenkalti/backoff/v4"
    "github.com/jackc/pgerrcode"
    "github.com/jackc/pgx/v5/pgconn"
    _ "github.com/jackc/pgx/v5/stdlib"
    "github.com/pressly/goose/v3"
)

// Неиспользуемая переменная для проверки реализации интерфейса хранилища БД репозиторием
var _ interfaces.Repository = (*DBRepository)(nil)

//go:embed migrations/*.sql
var embedMigrations embed.FS

// conflictError служит для возврата из retry операции специфичной ошибки конфликта
type conflictError struct {
    url queries.Url // URL, который уже присутствует в БД
    err error       // Ошибка конфликта данных - её нельзя явно возвращать из retry операции
}

// DBRepository является БД хранилищем URL.
type DBRepository struct {
    db *sql.DB          // Соединение с БД
    q  *queries.Queries // Подготовленные запросы
}

// NewDBRepository создает новое БД хранилище URL.
func NewDBRepository(databaseDSN string) (*DBRepository, error) {
    // Открываем новое соединение
    db, err := sql.Open("pgx", databaseDSN)
    if err != nil {
        slog.Error(
            "failed to open DB connection",
            slog.String("error", err.Error()),
        )
        return nil, err
    }

    // Устанавливаем параметры соединений
    db.SetMaxIdleConns(5)
    db.SetMaxOpenConns(5)
    db.SetConnMaxIdleTime(1 * time.Second)
    db.SetConnMaxLifetime(30 * time.Second)

    // Создаем запросы
    var q *queries.Queries

    // Сначала пробуем подготовить стейтменты запросов
    q, err = queries.Prepare(context.Background(), db)
    if err != nil {
        // Не получилось подготовить запросы, будем жить без подготовленных...
        q = queries.New(db)
    }

    // Создаем само БД хранилище
    dbRepository := &DBRepository{
        db: db,
        q:  q,
    }

    // Выполняем миграции
    if err := dbRepository.bootstrap(databaseDSN); err != nil {
        return nil, err
    }

    return dbRepository, nil
}

// bootstrap выполняет миграции БД.
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

// Store сохраняет слайс переданных URL в БД хранилище.
func (r *DBRepository) Store(urls []*model.URL) (*model.URL, error) {
    // Создаем контекст
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    // Начинаем транзакцию
    tx, err := r.db.Begin()
    if err != nil {
        return nil, err
    }
    // Откладываем откат транзакции - если всё будет ок, эффекта не будет
    defer func() { _ = tx.Rollback() }()

    // Получаем подготовленные запросы с ранее открытой транзакцией
    qtx := r.q.WithTx(tx)

    // Ошибка Postgres чтобы отловить ошибку конфликта данных
    var pgErr *pgconn.PgError

    // Обходим полученный слайс URL и записываем в БД с использованием retry операций
    for _, url := range urls {
        // Создаем функцию retry операции - возвращает отдельно структуру конфликта и ошибку
        f := func() (conflictError, error) {
            // conflictError - структура конфликта данных с конфликтным URL и ошибкой.
            // Ошибку конфликта нельзя возвращать из функции retry операции,
            // потому что по ней будут выполняться повторные вызовы.
            // Возможные варианты:
            // 1. Конфликта нет, ошибки нет - возвращаем два nil, все хорошо;
            // 2. Конфликта нет или есть, есть ошибка - возвращаем структуру конфликта и ошибку. Должны выполняться повторные вызовы операции;
            // 3. Конфликт есть, ошибки нет - возвращаем структуру конфликта и nil. Конфликт обрабатываем отдельно, повторных вызовов операции уже не должно быть.
            var ce conflictError

            // Пробуем выполнить запрос на сохранение URL
            // Если ранее с подготовкой было всё ок, выполняются подготовленные запросы в текущей транзакции
            _, errbo := qtx.StoreURL(ctx, queries.StoreURLParams{
                LongUrl: url.Long,
                BaseUrl: url.Base,
                UrlID:   url.ID,
            })
            // Проверяем ошибку на конфликт
            if errors.As(errbo, &pgErr) && pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
                // Это конфликт
                // При любой ошибке транзакцию надо откатывать
                errrb := tx.Rollback()
                if errrb != nil {
                    return ce, errrb
                }
                // Получаем данные конфликтного URL по оригинальному адресу при помощи retry операции
                urlByLong, errgbl := backoff.RetryWithData(func() (queries.Url, error) {
                    return r.q.GetURLByLong(ctx, url.Long)
                }, backoff.NewExponentialBackOff())
                // Проверяем ошибку получения конфликтного URL
                if errgbl != nil {
                    return ce, errgbl
                }
                // Возвращаем данные конфликтного URL и кастомную ошибку конфликта
                return conflictError{
                    url: urlByLong,
                    err: ErrConflict,
                }, nil
            }
            // Возврат из retry операции
            return ce, errbo
        }

        // Выполняем подготовленную операцию
        conflError, err := backoff.RetryWithData(f, backoff.NewExponentialBackOff())
        if err != nil {
            return nil, err
        }

        // Если был конфликт, возвращаем кастомную ошибку конфликта и данные конфликтного URL
        if errors.Is(conflError.err, ErrConflict) {
            return &model.URL{
                Long: conflError.url.LongUrl,
                Base: conflError.url.BaseUrl,
                ID:   conflError.url.UrlID}, conflError.err
        }
    }

    // Все было хорошо, ошибок нет, коммитим транзакцию и возвращаем
    return nil, tx.Commit()
}

// Get возвращает из БД хранилища данные URL по переданному идентификатору.
func (r *DBRepository) Get(id string) (*model.URL, error) {
    // Создаем контекст
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    // Получаем из БД данные URL при помощи retry операции
    url, err := backoff.RetryWithData(func() (queries.Url, error) {
        return r.q.GetURL(ctx, id)
    }, backoff.NewExponentialBackOff())
    if err != nil {
        return nil, err
    }

    // Возвращаем данные URL
    return &model.URL{
        Long: url.LongUrl,
        Base: url.BaseUrl,
        ID:   url.UrlID,
    }, nil
}

// Close закрывает соединение с БД.
func (r *DBRepository) Close() error {
    return r.db.Close()
}

// Check выполняет пинг БД.
func (r *DBRepository) Check() error {
    return r.db.Ping()
}
