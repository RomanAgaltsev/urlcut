// Пакет model содержит структуры данных, используемые приложением в работе.
package model

// Структуры, используемые приложением.
type (
	// URLDTO - структура тела запроса с оригинальным URL.
	URLDTO struct {
		URL string `json:"url"`
	}

	// ResultDTO - структура тела ответа с сокращенным URL.
	ResultDTO struct {
		Result string `json:"result"`
	}

	// IncomingBatchDTO - структура элемента батча тела запроса.
	IncomingBatchDTO struct {
		CorrelationID string `json:"correlation_id"`
		OriginalURL   string `json:"original_url"`
	}

	// OutgoingBatchDTO - структура элемента батча тела ответа.
	OutgoingBatchDTO struct {
		CorrelationID string `json:"correlation_id"`
		ShortURL      string `json:"short_url"`
	}

	// UserURLDTO - структура URL пользователя.
	UserURLDTO struct {
		ShortURL    string `json:"short_url"`
		OriginalURL string `json:"original_url"`
	}

	// ShortURLsDTO - структура тела запроса на удаление URL пользователя.
	ShortURLsDTO struct {
		IDs []string
	}
)
