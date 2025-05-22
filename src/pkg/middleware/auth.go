package middleware

import (
	"context"
	"net/http"
	"smart-contract-automation/src/internal/model"
	"smart-contract-automation/src/internal/service"
	"smart-contract-automation/src/pkg/common"
	"smart-contract-automation/src/pkg/logger"
	"strings"

	"go.uber.org/zap"
)

type contextKey string

const (
	UserContextKey contextKey = "user"
)

func AuthMiddleware(authService service.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				common.RespondWithError(w, http.StatusUnauthorized, "No authorization header")
				return
			}

			tokenParts := strings.Split(authHeader, " ")
			if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
				common.RespondWithError(w, http.StatusUnauthorized, "Invalid authorization header format")
				return
			}

			token, err := authService.ValidateToken(tokenParts[1])
			if err != nil {
				logger.Error("Invalid token", zap.Error(err))
				common.RespondWithError(w, http.StatusUnauthorized, "Invalid token")
				return
			}

			claims, ok := token.Claims.(*service.Claims)
			if !ok {
				common.RespondWithError(w, http.StatusUnauthorized, "Invalid token claims")
				return
			}

			ctx := context.WithValue(r.Context(), UserContextKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func RequireRole(roles ...model.Role) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := r.Context().Value(UserContextKey).(*service.Claims)
			if !ok {
				common.RespondWithError(w, http.StatusUnauthorized, "No user in context")
				return
			}

			hasRole := false
			for _, role := range roles {
				if claims.Role == role {
					hasRole = true
					break
				}
			}

			if !hasRole {
				common.RespondWithError(w, http.StatusForbidden, "Insufficient permissions")
				return

			}

			next.ServeHTTP(w, r)
		})
	}
}
