package routes

import (
	"github.com/Alireza7997/go_microservices/gateway/handlers"
	"github.com/Alireza7997/go_microservices/pkg/router"
)

func InitRoutes(router *router.Router) {
	router.Handle(`/api/auth/register/`, handlers.Register, "POST")
	router.Handle(`/api/auth/login/`, handlers.Login, "POST")

	router.Handle(`/api/chat/messages/`, handlers.PostMessage, "POST")
	router.Handle(`/api/chat/messages/{room}/`, handlers.GetMessages, "GET")

	router.Handle(`/api/greet/ping/`, handlers.Ping, "GET")
}
