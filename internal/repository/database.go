package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/RomanAgaltsev/urlcut/internal/database/queries"
	"github.com/RomanAgaltsev/urlcut/internal/interfaces"
	"github.com/RomanAgaltsev/urlcut/internal/model"

	"github.com/cenkalti/backoff/v4"
	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
)

// Неиспользуемая переменная для проверки реализации интерфейса хранилища БД репозиторием
var _ interfaces.Repository = (*DBRepository)(nil)

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
func NewDBRepository(db *sql.DB) (*DBRepository, error) {
	// Создаем запросы
	var q *queries.Queries

	// Сначала пробуем подготовить стейтменты запросов
	q, err := queries.Prepare(context.Background(), db)
	if err != nil {
		// Не получилось подготовить запросы, будем жить без подготовленных...
		q = queries.New(db)
	}

	// Создаем само БД хранилище
	dbRepository := &DBRepository{
		db: db,
		q:  q,
	}

	return dbRepository, nil
}

// Store сохраняет слайс переданных URL в БД хранилище.
func (r *DBRepository) Store(ctx context.Context, urls []*model.URL) (*model.URL, error) {
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
				Uid:     url.UID,
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
				ID:   conflError.url.UrlID,
				UID:  conflError.url.Uid}, conflError.err
		}
	}

	// Все было хорошо, ошибок нет, коммитим транзакцию и возвращаем
	return nil, tx.Commit()
}

// Get возвращает из БД хранилища данные URL по переданному идентификатору.
func (r *DBRepository) Get(ctx context.Context, id string) (*model.URL, error) {
	// Получаем из БД данные URL при помощи retry операции
	url, err := backoff.RetryWithData(func() (queries.Url, error) {
		return r.q.GetURL(ctx, id)
	}, backoff.NewExponentialBackOff())
	if err != nil {
		return nil, err
	}

	// Возвращаем данные URL
	return &model.URL{
		Long:    url.LongUrl,
		Base:    url.BaseUrl,
		ID:      url.UrlID,
		Deleted: url.IsDeleted,
	}, nil
}

func (r *DBRepository) GetUserURLs(ctx context.Context, uid uuid.UUID) ([]*model.URL, error) {
	// Получаем из БД URL по идентификатору пользователя при помощи retry операции
	urlsQuery, err := backoff.RetryWithData(func() ([]queries.Url, error) {
		return r.q.GetUserURLs(ctx, uid)
	}, backoff.NewExponentialBackOff())
	if err != nil {
		return nil, err
	}

	// Создаем слайс для возврата ссылок пользователя
	urls := make([]*model.URL, 0, len(urlsQuery))

	// Перекладываем URL из результата запроса в слайс
	for _, url := range urlsQuery {
		// Проверяем, не отменили ли контекст
		select {
		case <-ctx.Done():
			// Контекст отменили - возвращаем ничего и ошибку контекста
			return nil, ctx.Err()
		default:
			// Контекст не отменили - продолжаем перекладывать
			urls = append(urls, &model.URL{
				Long: url.LongUrl,
				Base: url.BaseUrl,
				ID:   url.UrlID,
				UID:  url.Uid,
			})
		}
	}

	return urls, nil
}

func (r *DBRepository) DeleteURLs(ctx context.Context, urls []*model.URL) error {
	// Начинаем транзакцию
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	// Откладываем откат транзакции - если всё будет ок, эффекта не будет
	defer func() { _ = tx.Rollback() }()

	// Получаем подготовленные запросы с ранее открытой транзакцией
	qtx := r.q.WithTx(tx)

	// Обходим полученный слайс URL и обновляем записи в БД с использованием retry операций
	for _, url := range urls {
		err = backoff.Retry(func() error {
			return qtx.DeleteURL(ctx, queries.DeleteURLParams{
				UrlID: url.ID,
				Uid:   url.UID,
			})
		}, backoff.NewExponentialBackOff())

		// Проверяем ошибку получения конфликтного URL
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// Close закрывает соединение с БД.
func (r *DBRepository) Close() error {
	err := r.q.Close()
	if err != nil {
		return err
	}
	return r.db.Close()
}
