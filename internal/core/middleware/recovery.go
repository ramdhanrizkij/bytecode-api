package middleware

import (
	"fmt"
	"runtime/debug"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	"github.com/ramdhanrizkij/bytecode-api/internal/shared/response"
)

// Recovery returns a Fiber middleware that catches panics, logs the panic value
// and full stack trace using the supplied zap logger, and converts the panic
// into a 500 Internal Server Error response so the server stays alive.
func Recovery(log *zap.Logger) fiber.Handler {
	return func(c *fiber.Ctx) (retErr error) {
		defer func() {
			if r := recover(); r != nil {
				// Capture the full goroutine stack trace.
				stack := debug.Stack()

				log.Error("panic recovered",
					zap.String("panic", fmt.Sprintf("%v", r)),
					zap.String("path", c.Path()),
					zap.String("method", c.Method()),
					zap.ByteString("stacktrace", stack),
				)

				retErr = response.Error(c, fiber.StatusInternalServerError, "internal server error")
			}
		}()

		return c.Next()
	}
}
