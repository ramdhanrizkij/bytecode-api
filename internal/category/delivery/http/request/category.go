package request

type UpsertCategoryRequest struct {
	Name        string  `json:"name" binding:"required,min=2,max=120"`
	Slug        string  `json:"slug" binding:"max=150"`
	Description *string `json:"description" binding:"omitempty,max=1000"`
	IsActive    bool    `json:"is_active"`
}
