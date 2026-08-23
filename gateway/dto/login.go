package dto

type LoginRequest struct {
	Username string `json:"username" g:"required"`
	Password string `json:"password" g:"max=32,required"`
}

type LoginResponse struct {
	UserID   int64  `json:"user_id"`
	Token    string `json:"token"`
	Username string `json:"username"`
}
