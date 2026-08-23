package calls

import (
	"microservice/greet_service/greet_pb"
)

// WithGreetService runs do with a greet service client.
func WithGreetService(do func(client greet_pb.GreetServiceClient)) {
	conn, err := dial("greet")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	do(greet_pb.NewGreetServiceClient(conn))
}
