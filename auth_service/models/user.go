package models

const UserTable = "users"

type User struct {
	Id               int64  `db:"id" json:"id" skipInsert:"+"`
	Username         string `db:"user_name" json:"user_name"`
	FirstName        string `db:"first_name" json:"first_name"`
	LastName         string `db:"last_name" json:"last_name"`
	DisplayName      string `json:"display_name" db:"display_name"`
	Email            string `db:"email" json:"email"`
	Password         string `db:"password" json:"-"`
	EmailConfirmed   bool   `db:"email_confirmed" json:"-" skipInsert:"+"`
	EmailConfirmCode string `db:"email_confirm_code" json:"-" skipInsert:"+"`
	JoinDate         int64  `db:"joined_date" json:"joined_date"`
}
