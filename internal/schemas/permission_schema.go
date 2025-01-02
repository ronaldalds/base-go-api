package schemas

type PermissionResponse struct {
	ID          uint   `json:"id"`
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ListPermissionRequest struct {
	Page  uint `query:"page" validate:"required,min=1"`
	Limit uint `query:"limit" validate:"required"`
}

type ListPermissionResponse struct {
	Page  uint                 `json:"page" validate:"required,min=1"`
	Limit uint                 `json:"limit" validate:"required"`
	Data  []PermissionResponse `json:"data"`
	Total uint                 `json:"total" validate:"required"`
}
