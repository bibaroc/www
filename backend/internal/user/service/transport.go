package service

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/bibaroc/www/backend/internal/user/pb"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	httptransport "github.com/go-kit/kit/transport/http"
)

func NewHTTPHandler(endpoints Set, logger log.Logger) http.Handler {
	options := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(errorEncoder),
		httptransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
	}

	m := http.NewServeMux()
	m.Handle("/create_user", httptransport.NewServer(
		endpoints.CreateUser,
		decodeHTTPCreateUserRequest,
		encodeHTTPGenericResponse,
		options...,
	))
	return m
}

var (
	_ pb.UserServiceServer = (*grpcServer)(nil)
)

type grpcServer struct {
	createUser grpctransport.Handler
}

func (s *grpcServer) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	return nil, nil
}

func NewGRPCServer(endpoints Set, logger log.Logger) pb.UserServiceServer {
	options := []grpctransport.ServerOption{
		grpctransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
	}
	return &grpcServer{
		createUser: grpctransport.NewServer(
			endpoints.CreateUser,
			decodeGRPCCreateUserRequest,
			encodeGRPCCreateUserResponse,
			options...,
		),
	}
}

func decodeHTTPCreateUserRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req pb.CreateUserRequest
	err := jsonDecodeValidate(r.Body, &req)
	return &req, err
}

func encodeHTTPGenericResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if f, ok := response.(endpoint.Failer); ok && f.Failed() != nil {
		errorEncoder(ctx, f.Failed(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func decodeGRPCCreateUserRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.CreateUserRequest)
	err := req.Validate()
	return req, err
}

func encodeGRPCCreateUserResponse(_ context.Context, response interface{}) (interface{}, error) {
	return response, nil
}
