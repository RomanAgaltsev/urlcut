package repository

import (
	"github.com/RomanAgaltsev/urlcut/internal/model"
)

type URLRepository interface {
	Store(url *model.URL) error
	Get(id string) (*model.URL, error)
	SaveState() error
	RestoreState() error
}
