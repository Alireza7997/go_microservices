package service

import (
	"context"
	"testing"

	"github.com/Alireza7997/go_microservices/chat_service/chat_pb"
)

func TestStoreInsertAndList(t *testing.T) {
	s := NewStore()

	m1 := s.Insert("general", "alice", "hi")
	m2 := s.Insert("general", "bob", "hello")

	if m1.ID >= m2.ID {
		t.Fatalf("expected increasing ids, got %d then %d", m1.ID, m2.ID)
	}
	if m1.SentAt.IsZero() || m2.SentAt.IsZero() {
		t.Fatal("expected timestamps to be set")
	}

	msgs := s.List("general")
	if len(msgs) != 2 {
		t.Fatalf("expected 2 messages, got %d", len(msgs))
	}
	if msgs[0].Body != "hi" || msgs[1].Body != "hello" {
		t.Fatalf("expected oldest-first order, got %v", msgs)
	}
}

func TestStoreRoomsAreIsolated(t *testing.T) {
	s := NewStore()

	s.Insert("a", "alice", "for a")
	s.Insert("b", "bob", "for b")

	if msgs := s.List("a"); len(msgs) != 1 || msgs[0].Body != "for a" {
		t.Fatalf("room isolation broken: %+v", msgs)
	}
	if msgs := s.List("unknown"); len(msgs) != 0 {
		t.Fatalf("expected empty room, got %+v", msgs)
	}
}

func TestStoreCapsMessagesPerRoom(t *testing.T) {
	s := NewStore()

	for i := 0; i < maxMessagesPerRoom+10; i++ {
		s.Insert("spam", "bot", "msg")
	}

	if got := len(s.List("spam")); got != maxMessagesPerRoom {
		t.Fatalf("expected cap of %d messages, got %d", maxMessagesPerRoom, got)
	}
}

func TestServiceRoundTrip(t *testing.T) {
	svc := &service{store: NewStore()}
	ctx := context.Background()

	if _, err := svc.SendMessage(ctx, &chat_pb.SendMessageRequest{Room: "dev", Username: "alice", Body: "hey"}); err != nil {
		t.Fatalf("SendMessage failed: %v", err)
	}

	res, err := svc.GetMessages(ctx, &chat_pb.GetMessagesRequest{Room: "dev"})
	if err != nil {
		t.Fatalf("GetMessages failed: %v", err)
	}
	if len(res.Messages) != 1 {
		t.Fatalf("expected 1 message, got %d", len(res.Messages))
	}
	m := res.Messages[0]
	if m.Room != "dev" || m.Username != "alice" || m.Body != "hey" || m.Id == 0 || m.SentAt == 0 {
		t.Fatalf("unexpected message: %+v", m)
	}
}
