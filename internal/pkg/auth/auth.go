// Пакет auth предоставляет инструменты для работы с авторизацией.
package auth

import (
	"net/http"

	"github.com/go-chi/jwtauth/v5"
	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

// UserIDKey используется для получения идентификатора пользователя из клеймов JWT-токена.
type UserIDKey string

const (
	// DefaultCookieName содержит имя куки по умолчанию.
	DefaultCookieName = "jwt"

	// DefaultCookiePath содержит путь куки по умолчанию.
	DefaultCookiePath = "/"

	// DefaultCookieMaxAge содержит возраст куки по умолчанию.
	DefaultCookieMaxAge = 3600

	// UserIDClaimName содержит имя ключа идентификатора пользователя в контексте.
	UserIDClaimName UserIDKey = "uid"
)

// NewJWTToken создает новый JWT токен.
func NewJWTToken(ja *jwtauth.JWTAuth) (token jwt.Token, tokenString string, err error) {
	return ja.Encode(map[string]interface{}{string(UserIDClaimName): uuid.New().String()})
}

// NewCookieWithDefaults создает новую куку со значениями по умолчанию и переданным в параметре значением.
func NewCookieWithDefaults(value string) *http.Cookie {
	return &http.Cookie{
		Name:     DefaultCookieName,
		Value:    value,
		Path:     DefaultCookiePath,
		MaxAge:   DefaultCookieMaxAge,
		SameSite: http.SameSiteDefaultMode,
	}
}
