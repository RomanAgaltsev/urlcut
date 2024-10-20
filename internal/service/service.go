package service

import (
	"github.com/RomanAgaltsev/urlcut/internal/model"
)

type URLService interface {
	Shorten(longURL string) (*model.URL, error)
	Expand(id string) (*model.URL, error)
}
