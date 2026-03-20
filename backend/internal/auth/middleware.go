package auth

import (
	"net/http"
)

// AuthMiddleware parses access token from httpOnly cookie and attaches claims to request context.
func AuthMiddleware(jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Allow CORS preflight.
			if r.Method == http.MethodOptions {
				next.ServeHTTP(w, r)
				return
			}

			// Public endpoints can skip this middleware.
			c, err := r.Cookie("access_token")
			if err != nil || c == nil {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			claims, err := ParseAccessToken(jwtSecret, c.Value)
			if err != nil {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			ctx := WithUserClaims(r.Context(), claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

