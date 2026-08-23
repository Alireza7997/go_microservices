package service

import (
	"context"

	"microservice/chat_service/chat_pb"
)

type service struct {
	chat_pb.UnimplementedChatServiceServer
	store *Store
}

func New() chat_pb.ChatServiceServer {
	return &service{store: NewStore()}
}

func (s *service) SendMessage(ctx context.Context, req *chat_pb.SendMessageRequest) (*chat_pb.SendMessageResponse, error) {
	msg := s.store.Insert(req.GetRoom(), req.GetUsername(), req.GetBody())

	return &chat_pb.SendMessageResponse{
		Message: &chat_pb.Message{
			Id:       msg.ID,
			Room:     msg.Room,
			Username: msg.Username,
			Body:     msg.Body,
			SentAt:   msg.SentAt.Unix(),
		},
	}, nil
}

func (s *service) GetMessages(ctx context.Context, req *chat_pb.GetMessagesRequest) (*chat_pb.GetMessagesResponse, error) {
	msgs := s.store.List(req.GetRoom())

	res := &chat_pb.GetMessagesResponse{}
	for _, m := range msgs {
		res.Messages = append(res.Messages, &chat_pb.Message{
			Id:       m.ID,
			Room:     m.Room,
			Username: m.Username,
			Body:     m.Body,
			SentAt:   m.SentAt.Unix(),
		})
	}

	return res, nil
}
