package model

import "fmt"

// URL - структура URL.
type URL struct {
	Long   string // Оригинальный URL
	Base   string // Базовый адрес сокращенного URL
	ID     string // Идентификатор оригинального URL в сокращенном
	CorrID string // Идентификатор для сопоставления элементов батча запроса и ответа
}

// Short возвращает строку сокращенного URL.
func (u *URL) Short() string {
	return fmt.Sprintf("%s/%s", u.Base, u.ID)
}
