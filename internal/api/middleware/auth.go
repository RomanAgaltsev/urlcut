// Пакет middleware реализует миддлваре для хендлеров запросов.
package middleware

import (
	"context"
	"net/http"

	"github.com/go-chi/jwtauth/v5"
	"github.com/lestrrat-go/jwx/v2/jwt"

	"github.com/RomanAgaltsev/urlcut/internal/pkg/auth"
)

// WithAuth возвращает хендлер, обернутый в миддлваре авторизации.
// Сначала выполняется поиск JWT токена с именем "jwt" в куки запроса.
// Далее, если токен в куки не найден, выполняется получение токена из заголовка "Authorization".
// Если токен не был найден, генерируется новый и в заголовки ответа добавляется кука с ним.
func WithAuth(ja *jwtauth.JWTAuth) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			// Переменные для токена
			var tokenString string
			var token jwt.Token
			var err error

			// Пробуем получить токен из куки
			tokenString = tokenFromCookie(r)
			if tokenString == "" {
				// Пробуем получить токен из заголовка Authorization
				tokenString = r.Header.Get("Authorization")
			}

			// Проверяем, удалось ли получить строку токена
			if tokenString == "" {
				// Получить строку токена не удалось, создаем новую
				token, tokenString, err = auth.NewJWTToken(ja)
				if err != nil {
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
					return
				}
			} else {
				// Декодируем токен
				token, _ = ja.Decode(tokenString)
				// Пробуем валидировать токен
				if err = jwt.Validate(token, ja.ValidateOptions()...); err != nil {
					token, tokenString, err = auth.NewJWTToken(ja)
					if err != nil {
						http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
						return
					}
				}
			}

			// Устанавливает куку в ответ
			http.SetCookie(w, auth.NewCookieWithDefaults(tokenString))
			w.Header().Set("Authorization", tokenString)

			// Валидацию прошли, получим утверждения
			claims := token.PrivateClaims()

			// uid пользователя передаем дальше через контекст
			ctx := context.WithValue(r.Context(), auth.UserIDClaimName, claims[string(auth.UserIDClaimName)])
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}

// WithID работает аналогично WithAuth, но, при отсутствии токена в куке или заголовке, не выдает новый,
// а возвращает статус Unauthorized - отказывает в авторизации запроса.
func WithID(ja *jwtauth.JWTAuth) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			// Пробуем получить токен из куки
			tokenString := tokenFromCookie(r)
			if tokenString == "" {
				// Пробуем получить токен из заголовка Authorization
				tokenString = r.Header.Get("Authorization")
			}

			// Если строки токена нет - нет авторизации
			if tokenString == "" {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			// Декодируем токен
			token, err := ja.Decode(tokenString)

			// Если токен не получилось декодировать - нет авторизации
			if err != nil {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			// Если токен не получилось валидировать - нет авторизации
			if err = jwt.Validate(token, ja.ValidateOptions()...); err != nil {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			// Валидацию прошли, получим утверждения
			claims := token.PrivateClaims()

			// Получаем идентификатор пользователя
			id, ok := claims[string(auth.UserIDClaimName)]
			// Если в утверждениях идентификатора нет или он пустой - нет авторизации
			if !ok || id == "" {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			// uid пользователя передаем дальше через контекст
			ctx := context.WithValue(r.Context(), auth.UserIDClaimName, id)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}

func tokenFromCookie(r *http.Request) string {
	for _, cookie := range r.Cookies() {
		if cookie.Name == "jwt" {
			return cookie.Value
		}
	}
	return ""
}
