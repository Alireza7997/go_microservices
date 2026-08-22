package errors

import (
	"errors"
	"testing"
)

func TestNewWithoutWrappedError(t *testing.T) {
	err := New(404, "page not found", "")

	if err == nil {
		t.Fatal("expected non-nil error")
	}

	want := "Code: 404, ErrMessage: page not found, ResponseMessage: , Err: <nil>"
	if err.Error() != want {
		t.Fatalf("unexpected Error() output:\n got: %s\nwant: %s", err.Error(), want)
	}
}

func TestNewWithWrappedError(t *testing.T) {
	wrapped := errors.New("db down")
	err := New(500, "db", "internal error", wrapped)

	var te *theError
	if !errors.As(err, &te) {
		t.Fatal("expected *theError")
	}
	if te.Err != wrapped {
		t.Fatalf("expected wrapped error to be kept, got %v", te.Err)
	}
	if te.StatusCode() != 500 {
		t.Fatalf("expected status code 500, got %d", te.StatusCode())
	}
	if te.Message() != "internal error" {
		t.Fatalf("expected message 'internal error', got %q", te.Message())
	}
}

func TestStatusErrorInterface(t *testing.T) {
	var _ StatusError = (*theError)(nil)

	err := New(403, "", "forbidden")
	se, ok := err.(StatusError)
	if !ok {
		t.Fatal("theError should implement StatusError")
	}
	if se.StatusCode() != 403 || se.Message() != "forbidden" {
		t.Fatalf("unexpected values: code=%d msg=%q", se.StatusCode(), se.Message())
	}
}
