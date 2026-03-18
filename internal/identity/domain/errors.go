package domain

import sharedErrors "github.com/ramdhanrizki/bytecode-api/internal/shared/errors"

var (
	ErrInvalidCredentials       = sharedErrors.Unauthorized("invalid email or password")
	ErrEmailAlreadyRegistered   = sharedErrors.Conflict("email is already registered")
	ErrEmailNotVerified         = sharedErrors.Forbidden("email is not verified")
	ErrUserNotFound             = sharedErrors.NotFound("user not found")
	ErrRoleNotFound             = sharedErrors.NotFound("role not found")
	ErrPermissionNotFound       = sharedErrors.NotFound("permission not found")
	ErrInvalidVerificationToken = sharedErrors.Validation("invalid verification token", map[string][]string{"token": {"token is invalid or expired"}})
	ErrInvalidRefreshToken      = sharedErrors.Unauthorized("invalid refresh token")
	ErrInactiveUser             = sharedErrors.Forbidden("user account is inactive")
	ErrUnauthenticated          = sharedErrors.Unauthorized("authentication required")
	ErrForbidden                = sharedErrors.Forbidden("forbidden")
)
