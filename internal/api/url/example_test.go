package url

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/go-chi/chi/v5"
	"github.com/go-resty/resty/v2"

	"github.com/RomanAgaltsev/urlcut/internal/config"
	"github.com/RomanAgaltsev/urlcut/internal/repository"
	"github.com/RomanAgaltsev/urlcut/internal/services"
)

func Example() {
	// Объявляем служебные константы
	const (
		// Путь к файловому хранилищу URL
		fileStoragePath = "storage.json"

		// Оригинальный URL, который необходимо сократить
		url = "https://practicum.yandex.ru/"
	)

	// Получаем конфигурацию приложения
	cfg, _ := config.Get()

	// Создаем in memory репозиторий
	repo := repository.NewInMemoryRepository(fileStoragePath)

	// Создаем сервис сокращения URL
	service, _ := services.NewShortener(repo, cfg)

	// Создаем роутер chi
	router := chi.NewRouter()

	// Создаем хендлеры запросов
	handlers := NewHandlers(service, cfg)

	// Настраиваем в роутере путь для хендлера
	router.Post("/", handlers.Shorten)

	// Создаем тестовый http-сервер
	httpSrv := httptest.NewServer(router)
	// Откладываем закрытие тестового http-сервера до выхода из функции
	defer httpSrv.Close()

	// Создаем клиент resty
	httpc := resty.New()

	// Создаем новый запрос и задаем его параметры
	req := httpc.R()
	req.Method = http.MethodPost
	req.URL = httpSrv.URL

	// Отправляем запрос и получаем ответ
	res, _ := req.
		SetHeader("Content-Type", ContentTypeText).
		SetBody(url).
		Send()

	// Если запрос успешно обработан, вернется статус Created
	if res.StatusCode() == http.StatusCreated {
		// Получаем сокращенный URL
		fmt.Println(string(res.Body()))
	}
}
