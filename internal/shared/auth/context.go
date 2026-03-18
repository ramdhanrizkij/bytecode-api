package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const ContextKeyAuthenticatedUser = "authenticated_user"

type AuthenticatedUser struct {
	ID          uuid.UUID
	Email       string
	Roles       []string
	Permissions []string
}

func SetAuthenticatedUser(c *gin.Context, user AuthenticatedUser) {
	c.Set(ContextKeyAuthenticatedUser, user)
}

func GetAuthenticatedUser(c *gin.Context) (AuthenticatedUser, bool) {
	value, ok := c.Get(ContextKeyAuthenticatedUser)
	if !ok {
		return AuthenticatedUser{}, false
	}

	user, ok := value.(AuthenticatedUser)
	return user, ok
}
