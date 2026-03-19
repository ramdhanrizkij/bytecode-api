package dto

import "time"

type ListInput struct {
	Page     int
	Limit    int
	Search   string
	Sort     string
	Order    string
	IsActive *bool
}

type PaginationMeta struct {
	Page       int
	Limit      int
	Total      int
	TotalPages int
}

type CategorySummary struct {
	ID          string
	Name        string
	Slug        string
	Description *string
	IsActive    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type CategoryListOutput struct {
	Categories []CategorySummary
	Meta       PaginationMeta
}

type CreateCategoryInput struct {
	Name        string
	Slug        string
	Description *string
	IsActive    bool
}

type UpdateCategoryInput struct {
	ID          string
	Name        string
	Slug        string
	Description *string
	IsActive    bool
}
