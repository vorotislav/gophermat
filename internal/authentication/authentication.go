package authentication

import (
	"errors"
	"fmt"
	"time"

	"gophermat/internal/models"

	"github.com/golang-jwt/jwt/v5"
)

const (
	Secret         = "super.secret.key"
	ExpiresAt      = time.Hour * 5
	LeewayDuration = 5
)

var (
	ErrTokenIsExpired       = errors.New("the token is expired")
	ErrParseToken           = errors.New("parse token error")
	ErrSignedToken          = errors.New("signed token")
	ErrUnknownSigningMethod = errors.New("unexpected signing method")
)

// Claims contains internal claims and user payload.
type Claims struct {
	jwt.RegisteredClaims
	models.TokenPayload
}

type Authenticator struct {
}

func NewAuthenticator() *Authenticator {
	return &Authenticator{}
}

func (a *Authenticator) GenerateToken(payload models.TokenPayload) (string, error) {
	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ExpiresAt)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		TokenPayload: payload,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(Secret))
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrSignedToken, err)
	}

	return ss, nil
}

func (a *Authenticator) ParseToken(tokenString string) (models.TokenPayload, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("%w: %v", ErrUnknownSigningMethod, token.Header["alg"])
		}

		return []byte(Secret), nil
	}, jwt.WithLeeway(LeewayDuration*time.Second))

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		if time.Now().After(claims.ExpiresAt.Time) {
			return models.TokenPayload{}, ErrTokenIsExpired
		}

		return claims.TokenPayload, nil
	}

	return models.TokenPayload{}, fmt.Errorf("%w: %w", ErrParseToken, err)
}
