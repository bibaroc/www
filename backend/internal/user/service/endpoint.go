package service

import (
	"context"
	"time"

	"github.com/bibaroc/www/backend/internal/user/pb"
	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/ratelimit"
	"github.com/sony/gobreaker"
	"golang.org/x/time/rate"
)

type Set struct {
	CreateUser endpoint.Endpoint
}

func NewEndpointSet(srv pb.UserServiceServer) Set {
	var createUser endpoint.Endpoint
	{
		createUser = MakeCreateUserEndpoint(srv)
		createUser = ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), 1))(createUser)
		createUser = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(createUser)
	}
	return Set{
		CreateUser: createUser,
	}
}

func MakeCreateUserEndpoint(srv pb.UserServiceServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*pb.CreateUserRequest)
		v, err := srv.CreateUser(ctx, req)
		if err != nil {
			return nil, err
		}
		return v, nil
	}
}
