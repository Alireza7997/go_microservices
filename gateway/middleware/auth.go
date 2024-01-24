package middleware

import (
	"context"
	"net/http"
	"service/auth/auth_protobuf"
	"service/gateway/calls"
	g "service/gateway/global"
	"service/pkg/errors"
)

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		token := r.Header.Get("Token")
		if token == "" {
			panic(errors.New(errors.InvalidStatus, errors.ReSignIn, "TokenNotSpecified", ""))
		}

		s := calls.NewAuthService()
		var user *auth_protobuf.User = nil
		s.Call(func(auth auth_protobuf.AuthServiceClient) {
			resGrpc, err := auth.Me(ctx, &auth_protobuf.MeRequest{AccessToken: token})
			if resGrpc != nil {
				s.Check(resGrpc.Error, err)
			} else {
				s.Check(nil, err)
			}

			user = resGrpc.User
		})

		ctx = context.WithValue(ctx, g.UserKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
