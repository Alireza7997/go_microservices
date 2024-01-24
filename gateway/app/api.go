package app

import (
	"fmt"
	"log"
	g "microservice/gateway/global"
	"microservice/gateway/routes"
	"microservice/pkg/router"
	"net/http"
)

func API() {
	mux := new(router.Router)
	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", g.CFG.Gateway.IP, g.CFG.Gateway.Port),
		Handler: mux,
	}
	// Server uses ServeHTTP(ResponseWriter, *Request) method

	g.Server = server

	// Router Settings
	routes.InitRoutes(mux)

	// Run App
	log.Panic(server.ListenAndServe().Error())
}
