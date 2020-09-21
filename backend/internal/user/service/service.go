package service

import (
	"context"

	"github.com/go-kit/kit/log"

	"github.com/bibaroc/www/backend/internal/user/pb"
)

var (
	_ pb.UserServiceServer = (*UserService)(nil)
)

func NewUserService(logger log.Logger) pb.UserServiceServer {
	return &UserService{logger: logger}
}

type UserService struct {
	logger log.Logger
}

func (us *UserService) CreateUser(context.Context, *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	_ = us.logger.Log("UserService", true)
	return nil, nil
}
