package interfaces

import (
	"github.com/RomanAgaltsev/urlcut/internal/model"
)

type URLShortExpander interface {
	Shorten(longURL string) (*model.URL, error)
	Expand(id string) (*model.URL, error)
	Check() error
}

type URLStoreGetter interface {
	Store(url *model.URL) error
	Get(id string) (*model.URL, error)
	Check() error
}

type StateSetGetter interface {
	SetState(state map[string]*model.URL) error
	GetState() map[string]*model.URL
}
