package middleware

import (
	"context"
	"github.com/sheinsviatoslav/gophermart/internal/auth"
	"net/http"
)

type contextKey string

const UserIDKey contextKey = "userID"

func WithAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := auth.ReadEncryptedCookie(r, "userID")
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
