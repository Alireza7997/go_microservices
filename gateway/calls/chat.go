package calls

import (
	"github.com/Alireza7997/go_microservices/chat_service/chat_pb"
)

// WithChatService runs do with a chat service client.
// Panics are recovered by the gateway's panic middleware
// and converted into JSON error responses.
func WithChatService(do func(client chat_pb.ChatServiceClient)) {
	conn, err := dial("chat")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	do(chat_pb.NewChatServiceClient(conn))
}
