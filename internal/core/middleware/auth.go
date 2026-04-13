package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"

	"github.com/ramdhanrizkij/bytecode-api/internal/shared/response"
	pkgjwt "github.com/ramdhanrizkij/bytecode-api/pkg/jwt"
)

const localKeyUser = "user"

// JWTAuth returns a Fiber middleware that validates the Bearer token in the
// Authorization header and stores the parsed *Claims in c.Locals("user").
// Requests without a valid token receive a 401 response.
func JWTAuth(secret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return response.Error(c, fiber.StatusUnauthorized, "missing authorization header")
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			return response.Error(c, fiber.StatusUnauthorized, "invalid authorization header format")
		}

		tokenStr := parts[1]
		claims, err := pkgjwt.ParseToken(tokenStr, secret)
		if err != nil {
			return response.Error(c, fiber.StatusUnauthorized, err.Error())
		}

		// Store parsed claims for downstream handlers and middleware.
		c.Locals(localKeyUser, claims)

		return c.Next()
	}
}

// GetCurrentUser retrieves the authenticated user's Claims from the Fiber
// context. Returns nil if the middleware was not applied or the type assertion
// fails (should never happen in normal flow).
func GetCurrentUser(c *fiber.Ctx) *pkgjwt.Claims {
	claims, _ := c.Locals(localKeyUser).(*pkgjwt.Claims)
	return claims
}
