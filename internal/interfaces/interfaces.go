package interfaces

import (
	"github.com/RomanAgaltsev/urlcut/internal/model"
)

type Service interface {
	Shorten(longURL string) (*model.URL, error)
	Expand(id string) (*model.URL, error)
	Close() error
	Check() error
}

type Repository interface {
	Store(url *model.URL) error
	Get(id string) (*model.URL, error)
	Close() error
	Check() error
}
