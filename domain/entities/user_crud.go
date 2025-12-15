package entities

type CreateUserRequest struct {
	Email    string  `json:"email"`
	Name     string  `json:"name"`
	Password string  `json:"password"`
	TypeIDs  []int64 `json:"type_ids"` // Array de IDs dos tipos (obrigatório pelo menos 1)
}

type UpdateUserRequest struct {
	Email   *string  `json:"email,omitempty"`
	Name    *string  `json:"name,omitempty"`
	TypeIDs *[]int64 `json:"type_ids,omitempty"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

type UserListResponse struct {
	ID        int64       `json:"id"`
	Email     string      `json:"email"`
	Name      string      `json:"name"`
	UserTypes []*UserType `json:"user_types"`
	CreatedAt string      `json:"created_at"`
}
