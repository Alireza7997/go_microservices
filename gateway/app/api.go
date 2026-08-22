package app

import (
	"fmt"
	"log/slog"
	"net/http"

	g "microservice/gateway/global"
	"microservice/gateway/middleware"
	"microservice/gateway/routes"
	"microservice/pkg/router"
)

func API() {
	mux := new(router.Router)
	mux.Middleware(middleware.Panic)
	mux.Middleware(middleware.Cors)
	mux.Middleware(middleware.Json)
	if g.CFG.MaxConcurrentRequests > 0 {
		mux.Middleware(middleware.ConcurrentLimiter(int(g.CFG.MaxConcurrentRequests)))
	}

	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", g.CFG.Gateway.IP, g.CFG.Gateway.Port),
		Handler: mux,
	}

	g.Server = server

	routes.InitRoutes(mux)

	slog.Info("gateway listening", "addr", server.Addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		slog.Error("gateway server failed", "err", err)
	}
}
