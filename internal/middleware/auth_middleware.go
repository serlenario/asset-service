package middleware

import (
	"context"
	"net/http"
	"strings"

	"asset-service/internal/service"
)

func Auth(svc *service.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h := r.Header.Get("Authorization")
			parts := strings.SplitN(h, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
				return
			}
			uid, err := svc.Validate(r.Context(), parts[1])
			if err != nil {
				http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), "uid", uid)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
