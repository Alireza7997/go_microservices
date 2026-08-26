package service

import (
	"sync"
	"time"
)

// Message is an in-memory chat message.
type Message struct {
	ID        int64
	Room      string
	Username  string
	Body      string
	SentAt    time.Time
}

const maxMessagesPerRoom = 100

// Store keeps chat messages in memory, grouped by room.
// State is intentionally not persisted — this service exists to
// demonstrate the microservice architecture.
type Store struct {
	mu       sync.RWMutex
	nextID   int64
	messages map[string][]Message
}

func NewStore() *Store {
	return &Store{messages: make(map[string][]Message)}
}

// Insert stores a message and returns it with id/timestamp filled in.
func (s *Store) Insert(room, username, body string) Message {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.nextID++
	msg := Message{
		ID:       s.nextID,
		Room:     room,
		Username: username,
		Body:     body,
		SentAt:   time.Now(),
	}

	s.messages[room] = append(s.messages[room], msg)
	if len(s.messages[room]) > maxMessagesPerRoom {
		s.messages[room] = s.messages[room][1:]
	}

	return msg
}

// List returns the messages of a room, oldest first.
func (s *Store) List(room string) []Message {
	s.mu.RLock()
	defer s.mu.RUnlock()

	msgs := s.messages[room]
	out := make([]Message, len(msgs))
	copy(out, msgs)

	return out
}
