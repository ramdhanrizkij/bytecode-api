package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims is the custom JWT payload embedded in every token issued by this service.
type Claims struct {
	UserID   string `json:"user_id"`
	Email    string `json:"email"`
	RoleName string `json:"role_name"`
	jwt.RegisteredClaims
}

// GenerateToken creates and signs a JWT token with the provided user information.
// The token expires after expiryHours hours.
func GenerateToken(userID, email, roleName, secret string, expiryHours int) (string, error) {
	claims := Claims{
		UserID:   userID,
		Email:    email,
		RoleName: roleName,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expiryHours) * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return signed, nil
}

// ParseToken validates and parses a JWT string, returning the embedded Claims.
// It returns descriptive errors for expired tokens, invalid signatures, and
// other malformed inputs.
func ParseToken(tokenString, secret string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&Claims{},
		func(token *jwt.Token) (interface{}, error) {
			// Enforce HMAC signing — reject any other algorithm (algorithm confusion attack).
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(secret), nil
		},
	)
	if err != nil {
		// Surface specific error types for cleaner handling upstream.
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, fmt.Errorf("token has expired")
		}
		if errors.Is(err, jwt.ErrTokenSignatureInvalid) {
			return nil, fmt.Errorf("invalid token signature")
		}
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}
