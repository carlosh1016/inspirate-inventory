// Package jwt issues and validates the HS256 access tokens used by the API.
// Refresh tokens are opaque and live in platform outside this package (see
// usecase/auth), not as JWTs.
package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const issuer = "inspirate-inventory"

// Claims are the custom claims embedded in every access token.
type Claims struct {
	UserID int64  `json:"uid"`
	Rol    string `json:"rol"`
	SedeID int64  `json:"sid"`
	jwt.RegisteredClaims
}

// ErrExpiredToken means the token's signature is valid but it has expired —
// callers should tell the client to hit /auth/refresh, not re-login.
var ErrExpiredToken = errors.New("access token expired")

// ErrInvalidToken covers any other parsing or validation failure.
var ErrInvalidToken = errors.New("invalid access token")

// Manager issues and validates access tokens.
type Manager interface {
	GenerateAccessToken(userID int64, rol string, sedeID int64) (token string, expiresAt time.Time, err error)
	ParseAccessToken(token string) (*Claims, error)
}

type hmacManager struct {
	secret []byte
	ttl    time.Duration
}

// New builds an HS256 Manager. ttl is the access token lifetime.
func New(secret string, ttl time.Duration) Manager {
	return &hmacManager{secret: []byte(secret), ttl: ttl}
}

func (m *hmacManager) GenerateAccessToken(userID int64, rol string, sedeID int64) (string, time.Time, error) {
	now := time.Now()
	expiresAt := now.Add(m.ttl)

	claims := Claims{
		UserID: userID,
		Rol:    rol,
		SedeID: sedeID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    issuer,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(m.secret)
	if err != nil {
		return "", time.Time{}, err
	}

	return signed, expiresAt, nil
}

func (m *hmacManager) ParseAccessToken(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return m.secret, nil
	}, jwt.WithIssuer(issuer))

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	if !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}
