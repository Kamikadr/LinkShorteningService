package authMiddleware

import (
	"context"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
	"strings"
)

type TokenParser interface {
	ParseToken(string) (string, error)
}

func New(logger *slog.Logger, parser TokenParser) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			const op = "middleware.authMiddleware.New()"
			logger.With("op", op, "request_id", middleware.GetReqID(r.Context()))
			token := r.Header.Get("Authorization")
			if token == "" {
				logger.Warn("token is empty")
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}
			token = strings.TrimPrefix(token, "Bearer ")
			tokenUser, err := parser.ParseToken(token)
			if err != nil {
				logger.Warn("token is not valid")
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}
			logger.Info("token is valid")
			ctx := context.WithValue(r.Context(), "user", tokenUser)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
