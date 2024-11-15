package middleware

import (
	"context"
	"net/http"

	"github.com/go-chi/jwtauth/v5"
	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

func WithAuth(ja *jwtauth.JWTAuth) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			// Пробуем получить токен
			tokenString := jwtauth.TokenFromCookie(r)
			token, err := ja.Decode(tokenString)
			if err != nil || token == nil {
				// Получить токен не удалось, выдаем куку
				token, tokenString, err = ja.Encode(map[string]interface{}{"uid": uuid.New().String()})
				if err != nil {
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
					return
				}
				http.SetCookie(w, &http.Cookie{
					Name:   "jwt",
					Value:  tokenString,
					Path:   "/",
					MaxAge: 300,
					Secure: true,
				})
			}

			// Пробуем валидировать токен
			if err = jwt.Validate(token, ja.ValidateOptions()...); err != nil {
				// Валидировать токен не удалось, выдаем куку
				token, tokenString, _ = ja.Encode(map[string]interface{}{"uid": uuid.New().String()})
				http.SetCookie(w, &http.Cookie{
					Name:   "jwt",
					Value:  tokenString,
					Path:   "/",
					MaxAge: 300,
					Secure: true,
				})
			}

			// Валидацию прошли, получим утверждения
			claims := token.PrivateClaims()

			// uid пользователя передаем дальше через контекст
			ctx := context.WithValue(r.Context(), "uid", claims["uid"])
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}
