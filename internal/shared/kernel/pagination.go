package kernel

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	DefaultPage  = 1
	DefaultLimit = 10
	MaxLimit     = 100
)

type PaginationQuery struct {
	Page   int
	Limit  int
	Search string
	Sort   string
	Order  string
}

func PaginationFromRequest(c *gin.Context) PaginationQuery {
	page := parsePositiveInt(c.DefaultQuery("page", "1"), DefaultPage)
	limit := parsePositiveInt(c.DefaultQuery("limit", "10"), DefaultLimit)
	if limit > MaxLimit {
		limit = MaxLimit
	}

	return PaginationQuery{
		Page:   page,
		Limit:  limit,
		Search: c.Query("search"),
		Sort:   c.Query("sort"),
		Order:  c.Query("order"),
	}
}

func TotalPages(total, limit int) int {
	if limit <= 0 {
		return 0
	}

	if total == 0 {
		return 0
	}

	pages := total / limit
	if total%limit != 0 {
		pages++
	}

	return pages
}

func parsePositiveInt(raw string, fallback int) int {
	value, err := strconv.Atoi(raw)
	if err != nil || value <= 0 {
		return fallback
	}

	return value
}
