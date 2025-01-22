package auth

import (
	"testing"

	"github.com/go-chi/jwtauth/v5"
)

func BenchmarkNewJWTToken(b *testing.B) {
	const secretKey = "secret key"
	tokenAuth := jwtauth.New("HS256", []byte(secretKey), nil)
	b.ResetTimer()
	for range b.N {
		_, _, _ = NewJWTToken(tokenAuth)
	}
}

func BenchmarkNewCookieWithDefaults(b *testing.B) {
	const secretKey = "secret key"
	tokenAuth := jwtauth.New("HS256", []byte(secretKey), nil)
	_, tokenString, _ := NewJWTToken(tokenAuth)
	b.ResetTimer()
	for range b.N {
		_ = NewCookieWithDefaults(tokenString)
	}
}
