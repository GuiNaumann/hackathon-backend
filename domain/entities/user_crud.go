package entities

type CreateUserRequest struct {
	Email    string  `json:"email"`
	Name     string  `json:"name"`
	Password string  `json:"password"`
	TypeIDs  []int64 `json:"type_ids"`
	SectorID *int64  `json:"sector_id,omitempty"`
}

type UpdateUserRequest struct {
	Email    *string  `json:"email,omitempty"`
	Name     *string  `json:"name,omitempty"`
	TypeIDs  *[]int64 `json:"type_ids,omitempty"`
	SectorID *int64   `json:"sector_id"` // REMOVER omitempty para permitir null
}

type UserListResponse struct {
	ID         int64       `json:"id"`
	Email      string      `json:"email"`
	Name       string      `json:"name"`
	UserTypes  []*UserType `json:"user_types"`
	SectorID   *int64      `json:"sector_id,omitempty"`
	SectorName string      `json:"sector_name,omitempty"`
	CreatedAt  string      `json:"created_at"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}
