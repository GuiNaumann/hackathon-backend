package entities

import "time"

type UserType struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type TypeUser struct {
	ID         int64     `json:"id"`
	UserID     int64     `json:"user_id"`
	UserTypeID int64     `json:"user_type_id"`
	CreatedAt  time.Time `json:"created_at"`
}

type UserTypePermission struct {
	ID         int64     `json:"id"`
	UserTypeID int64     `json:"user_type_id"`
	Endpoint   string    `json:"endpoint"`
	Method     string    `json:"method"`
	CreatedAt  time.Time `json:"created_at"`
}

type PersonalInformation struct {
	User      *User       `json:"user"`
	UserTypes []*UserType `json:"user_types"`
}
