package routes

import (
	"microservice/gateway/handlers"
	"microservice/pkg/router"
)

// Applies all necessary middlewares
// func middlewares(mux *router.Router) {

// }

func InitRoutes(router *router.Router) {
	router.Handle(`/api/auth/register/`, handlers.Register, "POST")
}
