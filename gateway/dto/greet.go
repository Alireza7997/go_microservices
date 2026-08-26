package dto

type PingResponse struct {
	Greeting   string `json:"greeting"`
	ServerTime int64  `json:"server_time"`
}
