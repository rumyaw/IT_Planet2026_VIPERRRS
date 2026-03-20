package auth

import (
	"context"
)

type UserClaims struct {
	UserID string
	Role   string
}

type ctxKey int

const userKey ctxKey = 0

func WithUserClaims(ctx context.Context, claims UserClaims) context.Context {
	return context.WithValue(ctx, userKey, claims)
}

func UserClaimsFromContext(ctx context.Context) (UserClaims, bool) {
	v := ctx.Value(userKey)
	if v == nil {
		return UserClaims{}, false
	}
	claims, ok := v.(UserClaims)
	return claims, ok
}

