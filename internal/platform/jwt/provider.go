package jwt

import (
	"fmt"
	"time"

	jwtv5 "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"github.com/ramdhanrizki/bytecode-api/configs"
	identityService "github.com/ramdhanrizki/bytecode-api/internal/identity/domain/service"
)

type Provider struct {
	secret    []byte
	issuer    string
	accessTTL time.Duration
}

type accessClaims struct {
	Email string   `json:"email"`
	Roles []string `json:"roles"`
	jwtv5.RegisteredClaims
}

func NewProvider(cfg configs.JWTConfig) *Provider {
	return &Provider{
		secret:    []byte(cfg.Secret),
		issuer:    cfg.Issuer,
		accessTTL: time.Duration(cfg.AccessTTLMinutes) * time.Minute,
	}
}

func (p *Provider) Generate(claims identityService.AccessTokenClaims) (identityService.AccessToken, error) {
	now := time.Now().UTC()
	expiresAt := now.Add(p.accessTTL)

	tokenClaims := accessClaims{
		Email: claims.Email,
		Roles: claims.Roles,
		RegisteredClaims: jwtv5.RegisteredClaims{
			Subject:   claims.UserID.String(),
			Issuer:    p.issuer,
			IssuedAt:  jwtv5.NewNumericDate(now),
			ExpiresAt: jwtv5.NewNumericDate(expiresAt),
		},
	}

	token := jwtv5.NewWithClaims(jwtv5.SigningMethodHS256, tokenClaims)
	signed, err := token.SignedString(p.secret)
	if err != nil {
		return identityService.AccessToken{}, err
	}

	return identityService.AccessToken{
		Token:     signed,
		ExpiresAt: expiresAt,
	}, nil
}

func (p *Provider) Parse(token string) (*identityService.ParsedAccessToken, error) {
	parsedToken, err := jwtv5.ParseWithClaims(token, &accessClaims{}, func(t *jwtv5.Token) (any, error) {
		if _, ok := t.Method.(*jwtv5.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return p.secret, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := parsedToken.Claims.(*accessClaims)
	if !ok || !parsedToken.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}

	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return nil, err
	}

	return &identityService.ParsedAccessToken{
		UserID: userID,
		Email:  claims.Email,
		Roles:  claims.Roles,
	}, nil
}
