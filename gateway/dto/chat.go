package dto

type PostMessageRequest struct {
	Room     string `json:"room" g:"required"`
	Username string `json:"username" g:"required"`
	Body     string `json:"body" g:"required"`
}

type Message struct {
	ID       int64  `json:"id"`
	Room     string `json:"room"`
	Username string `json:"username"`
	Body     string `json:"body"`
	SentAt   int64  `json:"sent_at"`
}
