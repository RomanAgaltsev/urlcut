package model

import (
	"fmt"
	"github.com/google/uuid"
)

// URL - структура URL.
type URL struct {
	Long   string    // Оригинальный URL
	Base   string    // Базовый адрес сокращенного URL
	ID     string    // Идентификатор оригинального URL в сокращенном
	CorrID string    // Идентификатор для сопоставления элементов батча запроса и ответа
	UID    uuid.UUID // Идентификатор пользователя
}

// Short возвращает строку сокращенного URL.
func (u *URL) Short() string {
	return fmt.Sprintf("%s/%s", u.Base, u.ID)
}
