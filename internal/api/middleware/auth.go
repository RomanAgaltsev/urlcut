package middleware

import (
	"context"
	"net/http"

	"github.com/go-chi/jwtauth/v5"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

func WithAuth(ja *jwtauth.JWTAuth) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			// Получим токен в виде строки
			tokenString := jwtauth.TokenFromCookie(r)
			if tokenString == "" {
				// Выдаем куку
				_, tokenString, _ = ja.Encode(map[string]interface{}{"uid": 1})
				http.SetCookie(w, &http.Cookie{
					Name:   "jwt",
					Value:  tokenString,
					MaxAge: 300,
				})
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			// Декодируем
			token, err := ja.Decode(tokenString)
			if err != nil || token == nil {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			// Валидируем
			if err := jwt.Validate(token, ja.ValidateOptions()...); err != nil {
				// Выдаем куку
				_, tokenString, _ = ja.Encode(map[string]interface{}{"uid": 1})
				http.SetCookie(w, &http.Cookie{
					Name:   "jwt",
					Value:  tokenString,
					MaxAge: 300,
				})
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
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
