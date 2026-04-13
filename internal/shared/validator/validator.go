package validator

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	"github.com/ramdhanrizkij/bytecode-api/internal/shared/response"
)

// validate is the package-level validator instance. It is initialised once and
// reused across all requests for performance.
var validate = validator.New()

// ValidationErrorDetail describes a single field-level validation failure.
type ValidationErrorDetail struct {
	Field   string `json:"field"`
	Tag     string `json:"tag"`
	Message string `json:"message"`
}

// ValidateStruct runs validation on the given struct and returns a slice of
// ValidationErrorDetail for every violated constraint. Returns nil if valid.
func ValidateStruct(s interface{}) []*ValidationErrorDetail {
	var errs validator.ValidationErrors
	if err := validate.Struct(s); err != nil {
		if !errors.As(err, &errs) {
			// Non-field validation error (e.g. invalid struct passed).
			return []*ValidationErrorDetail{
				{Field: "unknown", Tag: "unknown", Message: err.Error()},
			}
		}

		details := make([]*ValidationErrorDetail, 0, len(errs))
		for _, fe := range errs {
			details = append(details, &ValidationErrorDetail{
				Field:   fe.Field(),
				Tag:     fe.Tag(),
				Message: humanMessage(fe),
			})
		}
		return details
	}
	return nil
}

// ParseAndValidate parses the request JSON body into type T, runs struct
// validation, and — if any errors are found — immediately sends a 422 response
// and returns (nil, err). On success it returns (*T, nil).
func ParseAndValidate[T any](c *fiber.Ctx) (*T, error) {
	var body T

	if err := c.BodyParser(&body); err != nil {
		parseErr := fmt.Errorf("failed to parse request body: %w", err)
		_ = response.Error(c, fiber.StatusBadRequest, parseErr.Error())
		return nil, parseErr
	}

	if details := ValidateStruct(body); details != nil {
		_ = response.ValidationError(c, details)
		return nil, fmt.Errorf("validation failed")
	}

	return &body, nil
}

// humanMessage converts a validator.FieldError into a user-friendly Indonesian
// error message. Unknown tags fall back to the raw validator message.
func humanMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "Field ini wajib diisi"
	case "email":
		return "Format email tidak valid"
	case "min":
		return fmt.Sprintf("Minimal %s karakter", fe.Param())
	case "max":
		return fmt.Sprintf("Maksimal %s karakter", fe.Param())
	case "uuid", "uuid4":
		return "Format UUID tidak valid"
	default:
		return fe.Error()
	}
}
