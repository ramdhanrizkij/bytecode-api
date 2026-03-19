package request

type UpsertProductRequest struct {
	CategoryID  string  `json:"category_id" binding:"required,uuid"`
	Name        string  `json:"name" binding:"required,min=2,max=160"`
	Slug        string  `json:"slug" binding:"max=180"`
	Description *string `json:"description" binding:"omitempty,max=1000"`
	SKU         string  `json:"sku" binding:"required,min=2,max=120"`
	Price       int64   `json:"price" binding:"required,gte=0"`
	Stock       int     `json:"stock" binding:"required,gte=0"`
	IsActive    bool    `json:"is_active"`
}
