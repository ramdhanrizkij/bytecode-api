package response

import "github.com/gofiber/fiber/v2"

// Meta holds the status code and human-readable message included in every response.
type Meta struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Response is the standard envelope for all non-paginated API responses.
type Response struct {
	Meta Meta        `json:"meta"`
	Data interface{} `json:"data,omitempty"`
}

// PaginatedResponse wraps a list payload together with pagination metadata.
type PaginatedResponse struct {
	Meta       Meta            `json:"meta"`
	Data       interface{}     `json:"data"`
	Pagination *PaginationMeta `json:"pagination,omitempty"`
}

// PaginationMeta contains the information needed for a client to navigate pages.
type PaginationMeta struct {
	CurrentPage int   `json:"current_page"`
	PerPage     int   `json:"per_page"`
	TotalItems  int64 `json:"total_items"`
	TotalPages  int   `json:"total_pages"`
}

// Success sends a 200 OK response with the given data payload.
func Success(c *fiber.Ctx, message string, data interface{}) error {
	return c.Status(fiber.StatusOK).JSON(Response{
		Meta: Meta{Code: fiber.StatusOK, Message: message},
		Data: data,
	})
}

// Created sends a 201 Created response with the given data payload.
func Created(c *fiber.Ctx, message string, data interface{}) error {
	return c.Status(fiber.StatusCreated).JSON(Response{
		Meta: Meta{Code: fiber.StatusCreated, Message: message},
		Data: data,
	})
}

// SuccessWithPagination sends a 200 OK response that includes pagination metadata.
func SuccessWithPagination(c *fiber.Ctx, message string, data interface{}, pagination *PaginationMeta) error {
	return c.Status(fiber.StatusOK).JSON(PaginatedResponse{
		Meta:       Meta{Code: fiber.StatusOK, Message: message},
		Data:       data,
		Pagination: pagination,
	})
}

// Error sends a response with the given HTTP status code and error message.
func Error(c *fiber.Ctx, code int, message string) error {
	return c.Status(code).JSON(Response{
		Meta: Meta{Code: code, Message: message},
	})
}

// ValidationError sends a 422 Unprocessable Entity response with validation error details.
func ValidationError(c *fiber.Ctx, errors interface{}) error {
	return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
		"meta": Meta{
			Code:    fiber.StatusUnprocessableEntity,
			Message: "validation error",
		},
		"errors": errors,
	})
}
