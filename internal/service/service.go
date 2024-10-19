package service

import (
	"context"

	"github.com/RomanAgaltsev/urlcut/internal/model"
)

type URLService interface {
	Shorten(ctx context.Context, longURL string) (*model.URL, error)
	Expand(ctx context.Context, id string) (*model.URL, error)
}
