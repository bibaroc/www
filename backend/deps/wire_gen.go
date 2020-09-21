// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package deps

import (
	"github.com/bibaroc/www/backend/deps/providers"
	"github.com/bibaroc/www/backend/internal/user/pb"
	"github.com/bibaroc/www/backend/internal/user/service"
	"github.com/go-kit/kit/log"
	"github.com/google/wire"
)

// Injectors from wire.go:

func InjectLogger() log.Logger {
	logger := providers.NewKitLoggerSet()
	return logger
}

func InjectUserService() pb.UserServiceServer {
	logger := providers.NewKitLoggerSet()
	userServiceServer := service.NewUserService(logger)
	return userServiceServer
}

func InjectUserServiceMetrics() service.Metrics {
	metrics := service.MakeMetrics()
	return metrics
}

// wire.go:

var (
	loggerSet = wire.NewSet(providers.NewKitLoggerSet)
)
