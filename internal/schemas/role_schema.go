package schemas

type RoleResponse struct {
	ID          uint                 `json:"id"`
	Name        string               `json:"name"`
	Description string               `json:"description"`
	Permissions []PermissionResponse `json:"permissions"`
}

type ListRoleRequest struct {
	Page  uint `query:"page" validate:"required,min=1"`
	Limit uint `query:"limit" validate:"required"`
}

type ListRoleResponse struct {
	Page  uint           `json:"page"`
	Limit uint           `json:"limit"`
	Data  []RoleResponse `json:"data"`
	Total uint           `json:"total"`
}

type CreateRoleRequest struct {
	Name        string  `json:"name" validate:"required,min=3,max=100"`
	Description *string `json:"description"`
	Permissions []uint  `json:"permissions"`
}
