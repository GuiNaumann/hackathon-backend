package entities

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

type CreateUserRequest struct {
	Email    string  `json:"email"`
	Name     string  `json:"name"`
	Password string  `json:"password"`
	TypeIDs  []int64 `json:"type_ids"`
	SectorID *int64  `json:"sector_id,omitempty"` // NOVO
}

type UpdateUserRequest struct {
	Email    *string  `json:"email,omitempty"`
	Name     *string  `json:"name,omitempty"`
	TypeIDs  *[]int64 `json:"type_ids,omitempty"`
	SectorID *int64   `json:"sector_id,omitempty"` // NOVO
}

type UserListResponse struct {
	ID         int64       `json:"id"`
	Email      string      `json:"email"`
	Name       string      `json:"name"`
	UserTypes  []*UserType `json:"user_types"`
	SectorID   *int64      `json:"sector_id,omitempty"`   // NOVO
	SectorName string      `json:"sector_name,omitempty"` // NOVO
	CreatedAt  string      `json:"created_at"`
}
