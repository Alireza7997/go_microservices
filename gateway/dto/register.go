package dto

type RegisterRequest struct {
	Username        string `json:"username" g:"required"`
	Password        string `json:"password" g:"max=32,required"`
	PasswordConfirm string `json:"password_confirm" g:"max=32,required"`
	Email           string `json:"email" g:"max=32,required"`
}
