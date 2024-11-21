package middleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/RomanAgaltsev/urlcut/internal/pkg/auth"

	"github.com/go-chi/jwtauth/v5"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

func WithAuth(ja *jwtauth.JWTAuth) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			//			for name, values := range r.Header {
			//				// Loop over all values for the name.
			//				for _, value := range values {
			//					fmt.Println(name, value)
			//				}
			//			}
			// Пробуем получить токен
			tokenString := jwtauth.TokenFromCookie(r)
			//		fmt.Println(tokenString)
			token, err := ja.Decode(tokenString)
			if err != nil || token == nil {
				// Получить токен не удалось, выдаем куку
				token, tokenString, err = auth.NewJWTToken(ja)
				if err != nil {
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
					return
				}
				//http.SetCookie(w, auth.NewCookieWithDefaults(tokenString))
				w.Header().Set("Cookie", fmt.Sprintf("jwt=%s", tokenString))
			}

			// Пробуем валидировать токен
			if err = jwt.Validate(token, ja.ValidateOptions()...); err != nil {
				// Валидировать токен не удалось, выдаем куку
				token, tokenString, err = auth.NewJWTToken(ja)
				if err != nil {
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
					return
				}
				//http.SetCookie(w, auth.NewCookieWithDefaults(tokenString))
				w.Header().Set("Cookie", fmt.Sprintf("jwt=%s", tokenString))
			}

			// Валидацию прошли, получим утверждения
			claims := token.PrivateClaims()

			// uid пользователя передаем дальше через контекст
			ctx := context.WithValue(r.Context(), auth.UserIDClaimName, claims[string(auth.UserIDClaimName)])
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}
