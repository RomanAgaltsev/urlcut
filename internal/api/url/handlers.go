package url

import (
	"github.com/RomanAgaltsev/urlcut/internal/config"
	"github.com/RomanAgaltsev/urlcut/internal/interfaces"
)

const (
	ContentTypeJSON = "application/json"
	ContentTypeText = "text/plain; charset=utf-8"
)

type Handlers struct {
	shortener interfaces.Service
	cfg       *config.Config
}

func NewHandlers(shortener interfaces.Service, cfg *config.Config) *Handlers {
	return &Handlers{
		shortener: shortener,
		cfg:       cfg,
	}
}
