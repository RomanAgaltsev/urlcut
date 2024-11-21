package auth

import (
	"net/http"

	"github.com/go-chi/jwtauth/v5"
	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

type UserIDKey string

const (
	// DefaultCookieName содержит имя куки по умолчанию.
	DefaultCookieName = "jwt"

	// DefaultCookiePath содержит путь куки по умолчанию.
	DefaultCookiePath = "/"

	// DefaultCookieMaxAge содержит возраст куки по умолчанию.
	DefaultCookieMaxAge = 0

	// DefaultCookieSecure содержит признак защищенной куки по умолчанию.
	DefaultCookieSecure = true

	// UserIDClaimName содержит имя ключа идентификатора пользователя в контексте.
	UserIDClaimName UserIDKey = "uid"
)

func NewJWTToken(ja *jwtauth.JWTAuth) (token jwt.Token, tokenString string, err error) {
	return ja.Encode(map[string]interface{}{"uid": uuid.New().String()})
}

// NewCookieWithDefaults создает новую куку со значениями по умолчанию и переданным в параметре значением.
func NewCookieWithDefaults(value string) *http.Cookie {
	return &http.Cookie{
		Name:   DefaultCookieName,
		Value:  value,
		Path:   DefaultCookiePath,
		MaxAge: DefaultCookieMaxAge,
		Secure: DefaultCookieSecure,
	}
}
