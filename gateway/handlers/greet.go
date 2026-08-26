package handlers

import (
	"net/http"

	"github.com/Alireza7997/go_microservices/gateway/calls"
	"github.com/Alireza7997/go_microservices/gateway/dto"
	g "github.com/Alireza7997/go_microservices/gateway/global"
	"github.com/Alireza7997/go_microservices/gateway/handlers/utils"
	"github.com/Alireza7997/go_microservices/pkg/errors"

	greet_pb "github.com/Alireza7997/go_microservices/greet_service/greet_pb"
)

func ping(w http.ResponseWriter, r *http.Request) {
	var res *greet_pb.PingResponse
	calls.WithGreetService(func(client greet_pb.GreetServiceClient) {
		resp, err := client.Ping(r.Context(), &greet_pb.PingRequest{
			Name: r.URL.Query().Get("name"),
		})
		if err != nil {
			panic(errors.New(http.StatusServiceUnavailable, err.Error(), "greet service unavailable"))
		}
		res = resp
	})

	utils.WriteJSON(w, http.StatusOK, dto.PingResponse{
		Greeting:   res.Greeting,
		ServerTime: res.ServerTime,
	})
}

var Ping = g.Handler{Handler: ping}
