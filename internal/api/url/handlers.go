package url

import (
	"github.com/RomanAgaltsev/urlcut/internal/interfaces"
)

const (
	ContentTypeJSON = "application/json"
	ContentTypeText = "text/plain; charset=utf-8"
)

type Handlers struct {
	shortener interfaces.URLShortExpander
}

func New(shortener interfaces.URLShortExpander) *Handlers {
	return &Handlers{
		shortener: shortener,
	}
}
