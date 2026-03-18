package handler

import (
	"errors"
	"strings"

	"github.com/go-playground/validator/v10"

	sharedErrors "github.com/ramdhanrizki/bytecode-api/internal/shared/errors"
)

func validationError(err error) error {
	if err == nil {
		return sharedErrors.Unauthorized("authentication required")
	}

	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		fields := make(map[string][]string, len(validationErrors))
		for _, fieldErr := range validationErrors {
			field := toSnakeCase(fieldErr.Field())
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
	case "email":
		return field + " must be a valid email"
	case "min":
		return field + " must be at least " + param + " characters"
	case "max":
		return field + " must not exceed " + param + " characters"
	default:
		return field + " is invalid"
	}
}

func toSnakeCase(input string) string {
	var builder strings.Builder
	for i, r := range input {
		if i > 0 && r >= 'A' && r <= 'Z' {
			builder.WriteByte('_')
		}
		if r >= 'A' && r <= 'Z' {
			builder.WriteRune(r + 32)
			continue
		}
		builder.WriteRune(r)
	}
	return builder.String()
}
