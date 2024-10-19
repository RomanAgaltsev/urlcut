package repository

import (
	"context"

	"github.com/RomanAgaltsev/urlcut/internal/model"
)

type URLRepository interface {
	Store(ctx context.Context, url *model.URL) error
	Get(ctx context.Context, id string) (*model.URL, error)
}
