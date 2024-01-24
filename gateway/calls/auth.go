package calls

import (
	"errors"
	"fmt"
	"microservice/auth/auth_pb"
	g "microservice/gateway/global"
	"microservice/general"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type authService struct{}

var authServiceInstance = &authService{}

func NewAuthService() *authService {
	return authServiceInstance
}

func (authService) Call(do func(service auth_pb.AuthServiceClient)) {
	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", g.AuthService.IP, g.AuthService.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(errors.New("auth service not available"))
	}
	defer conn.Close()

	service := auth_pb.NewAuthServiceClient(conn)
	do(service)
}

func (authService) Check(err *general.Error, errCarrier error) {
	if err != nil {
		panic(errors.New(err.ErrMsg + "\n" + err.Message))
	}

	if errCarrier != nil && strings.Contains(errCarrier.Error(), "connection refused") {
		panic(errors.New("AuthServiceUnavailable"))
	}
	if errCarrier != nil {
		panic(errors.New("AuthServiceDidNotFinishProperly"))
	}

}
