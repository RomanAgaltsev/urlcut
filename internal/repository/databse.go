package repository

import (
	"context"
	"database/sql"
	"github.com/RomanAgaltsev/urlcut/internal/interfaces"
	"github.com/RomanAgaltsev/urlcut/internal/model"
)

var _ interfaces.URLStoreGetter = (*DBRepository)(nil)

type DBRepository struct {
	db *sql.DB
}

func NewDBRepository(db *sql.DB) *DBRepository {
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

func (r *DBRepository) Check() error {
	return r.db.PingContext(context.Background())
}
