package calls

import (
	"fmt"

	g "github.com/Alireza7997/go_microservices/gateway/global"
	"github.com/Alireza7997/go_microservices/pkg/errors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// dial opens an insecure gRPC connection to a microservice
// configured in env.yaml by its name.
func dial(name string) (*grpc.ClientConn, error) {
	svc, ok := g.CFG.Microservices[name]
	if !ok {
		return nil, errors.New(500, "config", "service not configured")
	}

	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", svc.IP, svc.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, errors.New(503, err.Error(), "service unavailable")
	}

	return conn, nil
}
