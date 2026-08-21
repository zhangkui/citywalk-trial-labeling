package middleware

import (
	"net/http"

	"citywalk/internal/platform/crypto"
)

func Auth(next http.HandlerFunc, parser func(*http.Request) (*crypto.Claims, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, err := parser(r)
		if err != nil { http.Error(w, "unauthorized", http.StatusUnauthorized); return }
		r = r.WithContext(crypto.WithClaims(r.Context(), claims))
		next(w, r)
	}
}

