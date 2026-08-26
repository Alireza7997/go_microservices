package utils

import (
	"encoding/json"
	"net/http"
)

// WriteJSON writes v as the response body with the given status code.
func WriteJSON(w http.ResponseWriter, code int, v interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		panic(err)
	}
}
