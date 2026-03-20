package auth

import (
	"net/http"

	"trumplin/internal/httputilx"
)

func RequireRoles(roles ...string) func(http.Handler) http.Handler {
	allowed := make(map[string]struct{}, len(roles))
	for _, r := range roles {
		allowed[r] = struct{}{}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := UserClaimsFromContext(r.Context())
			if !ok {
				httputilx.WriteError(w, http.StatusUnauthorized, "unauthorized")
				return
			}
			if _, exists := allowed[claims.Role]; !exists {
				httputilx.WriteError(w, http.StatusForbidden, "forbidden")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

