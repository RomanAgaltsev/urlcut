package url

import (
	"github.com/RomanAgaltsev/urlcut/internal/interfaces"
)

const (
	ContentTypeJSON = "application/json"
	ContentTypeText = "text/plain; charset=utf-8"
)

type Handlers struct {
	service interfaces.URLShortExpander
}

func New(service interfaces.URLShortExpander) *Handlers {
	return &Handlers{
		service: service,
	}
}
