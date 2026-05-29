package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	domainauth "github.com/masaru/strutfy-api/internal/domain/auth"
	"github.com/masaru/strutfy-api/internal/service"

	"github.com/julienschmidt/httprouter"
)

type contextKey string

const claimsContextKey contextKey = "auth_claims"

func Auth(tokenService *service.TokenService, next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			writeUnauthorized(w, "missing authorization header")
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			writeUnauthorized(w, "invalid authorization header")
			return
		}

		claims, err := tokenService.ParseAccessToken(parts[1])
		if err != nil {
			writeUnauthorized(w, "invalid or expired token")
			return
		}

		ctx := context.WithValue(r.Context(), claimsContextKey, claims)
		next(w, r.WithContext(ctx), ps)
	}
}

func RequireRole(role string, next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		claims, ok := ClaimsFromContext(r.Context())
		if !ok || !claims.HasRole(role) {
			writeForbidden(w, "insufficient role")
			return
		}

		next(w, r, ps)
	}
}

func RequirePermission(permission string, next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		claims, ok := ClaimsFromContext(r.Context())
		if !ok || !claims.HasPermission(permission) {
			writeForbidden(w, "insufficient permission")
			return
		}

		next(w, r, ps)
	}
}

func ClaimsFromContext(ctx context.Context) (*domainauth.Claims, bool) {
	claims, ok := ctx.Value(claimsContextKey).(*domainauth.Claims)
	return claims, ok
}

func CurrentUserID(ctx context.Context) string {
	claims, ok := ClaimsFromContext(ctx)
	if !ok {
		return ""
	}

	return claims.UserID
}

func CurrentOrgID(ctx context.Context) string {
	claims, ok := ClaimsFromContext(ctx)
	if !ok {
		return ""
	}

	return claims.OrgID
}

func writeUnauthorized(w http.ResponseWriter, message string) {
	writeJSONError(w, http.StatusUnauthorized, message)
}

func writeForbidden(w http.ResponseWriter, message string) {
	writeJSONError(w, http.StatusForbidden, message)
}

func writeJSONError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(map[string]string{
		"error": message,
	})
}