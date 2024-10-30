package url

import (
	"github.com/RomanAgaltsev/urlcut/internal/service"
)

const (
	ContentTypeJSON = "application/json"
	ContentTypeText = "text/plain; charset=utf-8"
)

type Handlers struct {
	service service.URLService
}

func New(service service.URLService) *Handlers {
	return &Handlers{
		service: service,
	}
}
