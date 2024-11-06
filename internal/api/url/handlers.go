package url

import (
	"github.com/RomanAgaltsev/urlcut/internal/interfaces"
)

const (
	ContentTypeJSON = "application/json"
	ContentTypeText = "text/plain; charset=utf-8"
)

type Handlers struct {
	shortener interfaces.Service
}

func New(shortener interfaces.Service) *Handlers {
	return &Handlers{
		shortener: shortener,
	}
}
