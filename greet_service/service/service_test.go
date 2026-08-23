package service

import (
	"context"
	"testing"

	"microservice/greet_service/greet_pb"
)

func TestPingWithAndWithoutName(t *testing.T) {
	svc := service{}

	res, err := svc.Ping(context.Background(), &greet_pb.PingRequest{Name: "alireza"})
	if err != nil {
		t.Fatalf("Ping failed: %v", err)
	}
	if res.Greeting != "hello, alireza!" || res.ServerTime == 0 {
		t.Fatalf("unexpected response: %+v", res)
	}

	res, err = svc.Ping(context.Background(), &greet_pb.PingRequest{})
	if err != nil {
		t.Fatalf("Ping failed: %v", err)
	}
	if res.Greeting != "hello, stranger!" {
		t.Fatalf("unexpected default greeting: %q", res.Greeting)
	}
}
