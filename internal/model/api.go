package model

type (
	// Request - структура тела запроса с оригинальным URL.
	Request struct {
		URL string `json:"url"`
	}

	// Response - структура тела ответа с сокращенным URL.
	Response struct {
		Result string `json:"result"`
	}

	// BatchRequest - структура элемента батча тела запроса.
	BatchRequest struct {
		CorrelationID string `json:"correlation_id"`
		OriginalURL   string `json:"original_url"`
	}

	// BatchResponse - структура элемента батча тела ответа.
	BatchResponse struct {
		CorrelationID string `json:"correlation_id"`
		ShortURL      string `json:"short_url"`
	}

	UserURL struct {
		ShortURL    string `json:"short_url"`
		OriginalURL string `json:"original_url"`
	}
)
