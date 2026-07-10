package handler

import "github.com/gofiber/fiber/v3"

// RegisterRoutes registers auth-related routes on the provided router group.
// Routes are intentionally unauthenticated — users need these to obtain a token.
func RegisterRoutes(router fiber.Router, handler *AuthHTTPHandler) {
	auth := router.Group("/auth")
	auth.Post("/register", handler.Register)
	auth.Post("/login", handler.Login)
	auth.Post("/refresh", handler.Refresh)
	auth.Post("/logout", handler.Logout)
}
