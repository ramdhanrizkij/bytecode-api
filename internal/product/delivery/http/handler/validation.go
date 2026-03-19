package handler

import (
	"errors"

	"github.com/go-playground/validator/v10"

	sharedErrors "github.com/ramdhanrizki/bytecode-api/internal/shared/errors"
	"github.com/ramdhanrizki/bytecode-api/internal/shared/textcase"
)

func validationError(err error) error {
	if err == nil {
		return sharedErrors.Unauthorized("authentication required")
	}

	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		fields := make(map[string][]string, len(validationErrors))
		for _, fieldErr := range validationErrors {
			field := textcase.ToSnakeCase(fieldErr.Field())
			fields[field] = append(fields[field], messageForTag(field, fieldErr.Tag(), fieldErr.Param()))
		}
		return sharedErrors.Validation("validation failed", fields)
	}

	return sharedErrors.Validation("validation failed", map[string][]string{"request": {err.Error()}})
}

func messageForTag(field, tag, param string) string {
	switch tag {
	case "required":
		return field + " is required"
	case "min":
		return field + " must be at least " + param + " characters"
	case "max":
		return field + " must not exceed " + param + " characters"
	case "uuid":
		return field + " must be a valid uuid"
	case "gte":
		return field + " must be at least " + param
	default:
		return field + " is invalid"
	}
}
