package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AccessClaims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func NewAccessToken(secret string, userID string, role string, ttlSeconds int) (string, error) {
	now := time.Now()
	claims := AccessClaims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(ttlSeconds) * time.Second)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func ParseAccessToken(secret string, tokenString string) (UserClaims, error) {
	if tokenString == "" {
		return UserClaims{}, errors.New("empty_token")
	}

	tok, err := jwt.ParseWithClaims(tokenString, &AccessClaims{}, func(t *jwt.Token) (any, error) {
		if t.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, errors.New("unexpected_signing_method")
		}
		return []byte(secret), nil
	})
	if err != nil {
		return UserClaims{}, err
	}

	claims, ok := tok.Claims.(*AccessClaims)
	if !ok || !tok.Valid {
		return UserClaims{}, errors.New("invalid_claims")
	}

	return UserClaims{
		UserID: claims.UserID,
		Role:   claims.Role,
	}, nil
}

