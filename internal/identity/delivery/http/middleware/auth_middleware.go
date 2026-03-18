package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	identityDomain "github.com/ramdhanrizki/bytecode-api/internal/identity/domain"
	identityService "github.com/ramdhanrizki/bytecode-api/internal/identity/domain/service"
	sharedAuth "github.com/ramdhanrizki/bytecode-api/internal/shared/auth"
	sharedResponse "github.com/ramdhanrizki/bytecode-api/internal/shared/response"
)

type AuthMiddleware struct {
	tokenProvider identityService.AccessTokenProvider
}

func NewAuthMiddleware(tokenProvider identityService.AccessTokenProvider) *AuthMiddleware {
	return &AuthMiddleware{tokenProvider: tokenProvider}
}

func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := bearerToken(c.GetHeader("Authorization"))
		if token == "" {
			sharedResponse.Error(c, identityDomain.ErrUnauthenticated)
			c.Abort()
			return
		}

		claims, err := m.tokenProvider.Parse(token)
		if err != nil {
			sharedResponse.Error(c, identityDomain.ErrUnauthenticated)
			c.Abort()
			return
		}

		sharedAuth.SetAuthenticatedUser(c, sharedAuth.AuthenticatedUser{
			ID:    claims.UserID,
			Email: claims.Email,
			Roles: claims.Roles,
		})
		c.Next()
	}
}

func bearerToken(header string) string {
	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
		return ""
	}
	return strings.TrimSpace(parts[1])
}

func Unauthorized(c *gin.Context) {
	sharedResponse.Error(c, identityDomain.ErrUnauthenticated)
	c.AbortWithStatus(http.StatusUnauthorized)
}
