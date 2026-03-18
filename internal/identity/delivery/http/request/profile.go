package request

type UpdateProfileRequest struct {
	FullName string `json:"full_name" binding:"required,min=3,max=120"`
}
