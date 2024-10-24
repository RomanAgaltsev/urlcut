package url

import (
	"github.com/RomanAgaltsev/urlcut/internal/service"
)

type Handlers struct {
	service service.URLService
}

func New(service service.URLService) *Handlers {
	return &Handlers{
		service: service,
	}
}
