package pagination

import (
	"fmt"
	"math"

	"github.com/gofiber/fiber/v2"
)

// PaginationQuery holds the parsed query parameters for paginated list requests.
type PaginationQuery struct {
	Page    int    `query:"page"`
	PerPage int    `query:"per_page"`
	Sort    string `query:"sort"`
	Order   string `query:"order"`
	Search  string `query:"search"`
}

// NewPaginationQuery parses and normalises pagination parameters from the Fiber
// request context, applying sensible defaults and enforcing constraints.
func NewPaginationQuery(c *fiber.Ctx) *PaginationQuery {
	pq := &PaginationQuery{}

	// Parse raw query values into the struct fields.
	_ = c.QueryParser(pq)

	// Defaults and bounds for Page.
	if pq.Page < 1 {
		pq.Page = 1
	}

	// Defaults and bounds for PerPage.
	if pq.PerPage < 1 {
		pq.PerPage = 10
	}
	if pq.PerPage > 100 {
		pq.PerPage = 100
	}

	// Default sort column.
	if pq.Sort == "" {
		pq.Sort = "created_at"
	}

	// Default and validate sort direction.
	if pq.Order != "asc" && pq.Order != "desc" {
		pq.Order = "desc"
	}

	return pq
}

// GetOffset returns the number of records to skip for the current page.
func (p *PaginationQuery) GetOffset() int {
	return (p.Page - 1) * p.PerPage
}

// GetLimit returns the maximum number of records to return per page.
func (p *PaginationQuery) GetLimit() int {
	return p.PerPage
}

// GetSort returns the GORM-compatible ORDER BY clause string,
// e.g. "created_at desc".
func (p *PaginationQuery) GetSort() string {
	return fmt.Sprintf("%s %s", p.Sort, p.Order)
}

// CalculateTotalPages returns the total number of pages given the total item
// count and per-page size. Always returns at least 1.
func CalculateTotalPages(totalItems int64, perPage int) int {
	if perPage <= 0 || totalItems <= 0 {
		return 1
	}
	return int(math.Ceil(float64(totalItems) / float64(perPage)))
}
