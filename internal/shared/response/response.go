package response

import (
	"net/http"

	"github.com/gin-gonic/gin"

	sharedErrors "github.com/ramdhanrizki/bytecode-api/internal/shared/errors"
)

type Meta struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

type body struct {
	Success bool                `json:"success"`
	Message string              `json:"message"`
	Data    any                 `json:"data,omitempty"`
	Errors  map[string][]string `json:"errors,omitempty"`
	Meta    *Meta               `json:"meta,omitempty"`
}

func Success(c *gin.Context, status int, message string, data any) {
	c.JSON(status, body{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func Paginated(c *gin.Context, message string, data any, meta Meta) {
	c.JSON(http.StatusOK, body{
		Success: true,
		Message: message,
		Data:    data,
		Meta:    &meta,
	})
}

func Error(c *gin.Context, err error) {
	if appErr, ok := err.(*sharedErrors.AppError); ok {
		c.JSON(appErr.StatusCode, body{
			Success: false,
			Message: appErr.Message,
			Errors:  appErr.Fields,
		})
		return
	}

	c.JSON(http.StatusInternalServerError, body{
		Success: false,
		Message: "internal server error",
	})
}
