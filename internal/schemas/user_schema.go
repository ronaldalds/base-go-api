package schemas

type CreateUser struct {
	ID          uint   `json:"id"`
	FirstName   string `json:"firstName" validate:"required,min=1,max=50"`
	LastName    string `json:"lastName" validate:"omitempty,max=50"`
	Username    string `json:"username" validate:"required,min=3,max=50"`
	Email       string `json:"email" validate:"required,email"`
	Password    string `json:"password" validate:"required,min=6"`
	Active      bool   `json:"active"`
	IsSuperUser bool   `json:"isSuperUser"`
	Roles       []uint `json:"roles"`
	Phone1      string `json:"phone1" validate:"required,e164"`
	Phone2      string `json:"phone2" validate:"omitempty,e164"`
}

type UpdateUser struct {
	ID          uint   `json:"id"`
	FirstName   string `json:"firstName" validate:"omitempty,min=3,max=50"`
	LastName    string `json:"lastName" validate:"omitempty,min=3,max=50"`
	Username    string `json:"username" validate:"omitempty,min=3,max=50"`
	Email       string `json:"email" validate:"omitempty,email"`
	Active      bool   `json:"active"`
	IsSuperUser bool   `json:"isSuperUser"`
	Roles       []uint `json:"roles"`
	Phone1      string `json:"phone1" validate:"omitempty,e164"`
	Phone2      string `json:"phone2" validate:"omitempty,e164"`
}

type UserResponse struct {
	ID          uint     `json:"id"`
	FirstName   string   `json:"firstName"`
	LastName    string   `json:"lastName"`
	Username    string   `json:"username"`
	Email       string   `json:"email"`
	Active      bool     `json:"active"`
	IsSuperUser bool     `json:"isSuperUser"`
	RoleNames   []string `json:"roleNames"`
	Phone1      string   `json:"phone1"`
	Phone2      string   `json:"phone2"`
}

type ListUsersRequest struct {
	Page  uint `query:"page" validate:"required,min=1"`
	Limit uint `query:"limit" validate:"required"`
}

type ListUsersResponse struct {
	Page  uint           `json:"page" validate:"required,min=1"`
	Limit uint           `json:"limit" validate:"required"`
	Data  []UserResponse `json:"data"`
	Total uint           `json:"total" validate:"required"`
}
