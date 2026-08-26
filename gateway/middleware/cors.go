package middleware

import (
	"net/http"
	"strings"

	g "github.com/Alireza7997/go_microservices/gateway/global"
)

var allowMethods = "GET, POST, PUT, PATCH, DELETE, OPTIONS"
var allowHeaders = "Origin, Content-Length, Content-Type"

func Cors(next http.Handler) http.Handler {
	allowedOrigins := strings.Split(g.CFG.AllowOrigins, ",")
	for i := range allowedOrigins {
		allowedOrigins[i] = strings.TrimSpace(allowedOrigins[i])
	}

	if g.CFG.AllowHeaders != "" {
		allowHeaders += ", " + g.CFG.AllowHeaders
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin == "" {
			next.ServeHTTP(w, r)
			return
		}

		host := r.Host
		if origin == "http://"+host || origin == "https://"+host {
			next.ServeHTTP(w, r)
			return
		}

		found := false
		for _, allowed := range allowedOrigins {
			if allowed == "*" || allowed == origin {
				found = true
				break
			}
		}

		if !found && len(allowedOrigins) != 0 {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		w.Header().Set("Vary", "Origin")
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", allowMethods)
		w.Header().Set("Access-Control-Allow-Headers", allowHeaders)

		if strings.ToUpper(r.Method) == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
