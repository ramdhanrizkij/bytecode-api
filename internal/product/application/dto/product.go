package dto

import "time"

type ListInput struct {
	Page       int
	Limit      int
	Search     string
	Sort       string
	Order      string
	CategoryID *string
	IsActive   *bool
}

type PaginationMeta struct {
	Page       int
	Limit      int
	Total      int
	TotalPages int
}

type ProductSummary struct {
	ID           string
	CategoryID   string
	CategoryName string
	Name         string
	Slug         string
	Description  *string
	SKU          string
	Price        int64
	Stock        int
	IsActive     bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type ProductListOutput struct {
	Products []ProductSummary
	Meta     PaginationMeta
}

type CreateProductInput struct {
	CategoryID  string
	Name        string
	Slug        string
	Description *string
	SKU         string
	Price       int64
	Stock       int
	IsActive    bool
}

type UpdateProductInput struct {
	ID          string
	CategoryID  string
	Name        string
	Slug        string
	Description *string
	SKU         string
	Price       int64
	Stock       int
	IsActive    bool
}
