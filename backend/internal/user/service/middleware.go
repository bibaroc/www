package service

import (
	"context"
	"fmt"
	"time"

	"github.com/bibaroc/www/backend/internal/user/pb"
	"github.com/go-kit/kit/log"
)

var (
	_ pb.UserServiceServer = (*loggingMiddleware)(nil)
	_ pb.UserServiceServer = (*instrumentingMiddleware)(nil)
)

func LoggingMiddleware(
	logger log.Logger,
) func(next pb.UserServiceServer) pb.UserServiceServer {
	return func(next pb.UserServiceServer) pb.UserServiceServer {
		return loggingMiddleware{logger, next}
	}
}

type loggingMiddleware struct {
	logger log.Logger
	next   pb.UserServiceServer
}

func (mw loggingMiddleware) CreateUser(
	ctx context.Context,
	req *pb.CreateUserRequest,
) (res *pb.CreateUserResponse, err error) {
	defer func() {
		_ = mw.logger.Log(
			"method", "CreateUser",
			"username", req.Username,
			"email", req.Email,
			"err", err)
	}()
	res, err = mw.next.CreateUser(ctx, req)
	return
}

type instrumentingMiddleware struct {
	Metrics
	next pb.UserServiceServer
}

func InstrumentingMiddleware(
	metrics Metrics,
) func(next pb.UserServiceServer) pb.UserServiceServer {
	return func(next pb.UserServiceServer) pb.UserServiceServer {
		return instrumentingMiddleware{
			Metrics: metrics,
			next:    next,
		}
	}
}
func (mw instrumentingMiddleware) CreateUser(
	ctx context.Context,
	req *pb.CreateUserRequest,
) (res *pb.CreateUserResponse, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "CreateUser", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	res, err = mw.next.CreateUser(ctx, req)
	return
}
