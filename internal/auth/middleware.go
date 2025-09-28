package auth

import (
	"context"
	"net/http"

	"jwt-golang/internal/common"
)

type ContextKey string

const (
	userContextKey   ContextKey = "context"
	userIdContextKey ContextKey = "userId"
)

// should just get the JWT data and bind it to the contextd
func JwtMiddleware(authManager AuthManager, app *common.Application) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authorizationHeader := r.Header.Get("Authorization")

			// validate token format and existence
			err := authManager.ValidateAccessToken(authorizationHeader)
			if err != nil {
				app.Logger.Debug("token validation failed", "error", err.Error())
				next.ServeHTTP(w, r)
				return
			}

			// extract token string from header
			tokenString, err := authManager.ExtractTokenFromHeader(authorizationHeader)
			if err != nil || tokenString == "" {
				app.Logger.Debug("no token found in header")
				next.ServeHTTP(w, r)
				return
			}

			// get claims from token
			claims, err := authManager.ExtractClaims(tokenString)
			if err != nil {
				app.Logger.Debug("failed to extract claims", "error", err.Error())
				next.ServeHTTP(w, r)
				return
			}

			// create context with user data
			ctx := context.WithValue(r.Context(), userContextKey, claims)
			ctx = context.WithValue(ctx, userIdContextKey, claims.Subject)

			app.Logger.Debug("user authenticated", "userId", claims.Subject, "email", claims.Email)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func RequireAuthentication(AuthManager AuthManager, app *common.Application) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user := r.Context().Value(userContextKey)
			if user == nil {
				app.ErrorResponse(w, r, http.StatusUnauthorized, "authentication required")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// helper functions to get user data from context
func GetUserFromContext(ctx context.Context) (*CustomClaims, bool) {
	user, ok := ctx.Value(userContextKey).(*CustomClaims)
	return user, ok
}

func GetUserIdFromContext(ctx context.Context) (string, bool) {
	userId, ok := ctx.Value(userIdContextKey).(string)
	return userId, ok
}
