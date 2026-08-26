package middleware

import (
	"encoding/json"
	"log/slog"
	"net/http"

	pkgErrors "github.com/Alireza7997/go_microservices/pkg/errors"
)

type ErrorResponse struct {
	Code    int64  `json:"code"`
	Message string `json:"message"`
}

func Panic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			recovered := recover()
			if recovered == nil {
				return
			}

			if statusErr, ok := recovered.(pkgErrors.StatusError); ok {
				code := statusErr.StatusCode()
				if code < 100 || code > 599 {
					code = http.StatusInternalServerError
				}
				writeError(w, int(code), ErrorResponse{Code: code, Message: statusErr.Message()})
				return
			}

			if err, ok := recovered.(error); ok {
				slog.Error("request panicked", "err", err)
			} else {
				slog.Error("request panicked", "err", recovered)
			}

			writeError(w, http.StatusInternalServerError, ErrorResponse{
				Code:    http.StatusInternalServerError,
				Message: "internal server error",
			})
		}()
		next.ServeHTTP(w, r)
	})
}

func writeError(w http.ResponseWriter, code int, res ErrorResponse) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(res)
}
