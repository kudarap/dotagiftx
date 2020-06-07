package jwt

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// SigKey JWT signature key and you need to set your own!
var SigKey = "mal3d1ct10n"

// Claims represents JWT payload.
type Claims struct {
	UserID string
	Level  string
	jwt.StandardClaims
}

// New creates new JWT with claims payload.
func New(userID, level string, expiration time.Time) (token string, err error) {
	c := Claims{UserID: userID, Level: level}
	c.ExpiresAt = expiration.Unix()
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	return t.SignedString([]byte(SigKey))
}

// Parse validates and extract claims from token.
func Parse(token string) (*Claims, error) {
	t, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(SigKey), nil
	})

	if t == nil || !t.Valid {
		return nil, fmt.Errorf("token is not valid: %s", err)
	}

	tc, ok := t.Claims.(*Claims)
	if !ok || tc == nil {
		return nil, errors.New("token have empty claims")
	}

	return tc, nil
}

// ParseFromHeader parses Authorization bearer token as JWT.
func ParseFromHeader(h http.Header) (*Claims, error) {
	// Get access header.
	header := strings.TrimSpace(h.Get("Authorization"))
	if header == "" {
		return nil, errors.New("empty access header")
	}

	// Should contain only Bearer and Token.
	parts := strings.Split(header, " ")
	if len(parts) != 2 {
		return nil, errors.New("invalid access header")
	}

	// Check for bearer existence.
	if strings.ToUpper(parts[0]) != "BEARER" {
		return nil, errors.New("no bearer on access header")
	}

	if strings.TrimSpace(parts[1]) == "" {
		return nil, errors.New("empty bearer token")
	}

	return Parse(parts[1])
}
