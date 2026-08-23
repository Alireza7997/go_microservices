package handlers

import (
	"net/http"

	"microservice/chat_service/chat_pb"
	"microservice/gateway/calls"
	"microservice/gateway/dto"
	g "microservice/gateway/global"
	"microservice/gateway/handlers/utils"
	"microservice/pkg/errors"
	"microservice/pkg/router"
)

func postMessage(w http.ResponseWriter, r *http.Request) {
	req := &dto.PostMessageRequest{}

	utils.ParseBody(r.Body, req)
	if req.Room == "" || req.Username == "" || req.Body == "" {
		panic(errors.New(http.StatusBadRequest, "missing field", "room, username and body are required"))
	}

	var res *chat_pb.SendMessageResponse
	calls.WithChatService(func(client chat_pb.ChatServiceClient) {
		resp, err := client.SendMessage(r.Context(), &chat_pb.SendMessageRequest{
			Room:     req.Room,
			Username: req.Username,
			Body:     req.Body,
		})
		if err != nil {
			panic(errors.New(http.StatusServiceUnavailable, err.Error(), "chat service unavailable"))
		}
		res = resp
	})

	utils.WriteJSON(w, http.StatusCreated, dto.Message{
		ID:       res.Message.Id,
		Room:     res.Message.Room,
		Username: res.Message.Username,
		Body:     res.Message.Body,
		SentAt:   res.Message.SentAt,
	})
}

func getMessages(w http.ResponseWriter, r *http.Request) {
	room, _ := r.Context().Value(router.ReturnContextKey("room")).(string)

	var res *chat_pb.GetMessagesResponse
	calls.WithChatService(func(client chat_pb.ChatServiceClient) {
		resp, err := client.GetMessages(r.Context(), &chat_pb.GetMessagesRequest{Room: room})
		if err != nil {
			panic(errors.New(http.StatusServiceUnavailable, err.Error(), "chat service unavailable"))
		}
		res = resp
	})

	msgs := make([]dto.Message, 0, len(res.Messages))
	for _, m := range res.Messages {
		msgs = append(msgs, dto.Message{
			ID:       m.Id,
			Room:     m.Room,
			Username: m.Username,
			Body:     m.Body,
			SentAt:   m.SentAt,
		})
	}

	utils.WriteJSON(w, http.StatusOK, msgs)
}

var PostMessage = g.Handler{Handler: postMessage}
var GetMessages = g.Handler{Handler: getMessages}
