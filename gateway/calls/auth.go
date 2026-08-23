package calls

import (
	"microservice/auth_service/auth_pb"
)

// WithAuthService runs do with an auth service client.
func WithAuthService(do func(client auth_pb.AuthServiceClient)) {
	conn, err := dial("auth")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	do(auth_pb.NewAuthServiceClient(conn))
}
