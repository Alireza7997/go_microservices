package models

import (
	"time"
)

type User struct {
	Id          uint64    `json:"id" db:"id" goqu:"skipinsert"`
	Username    string    `json:"username" db:"username"`
	Email       string    `json:"email" db:"email"`
	Password    string    `json:"password" db:"password"`
	DisplayName string    `json:"display_name" db:"display_name"`
	CreatedAt   time.Time `db:"created_at" goqu:"skipupdate"`
	UpdatedAt   time.Time `db:"updated_at"`
}
