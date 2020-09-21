//+build wireinject

package deps

import (
	"github.com/go-kit/kit/log"

	"github.com/bibaroc/www/backend/deps/providers"
	"github.com/bibaroc/www/backend/internal/user/pb"
	"github.com/bibaroc/www/backend/internal/user/service"
	"github.com/google/wire"
)

var (
	loggerSet = wire.NewSet(
		providers.NewKitLoggerSet,
	)
)

func InjectLogger() log.Logger {
	wire.Build(
		loggerSet,
	)
	return log.Logger(nil)
}

func InjectUserService() pb.UserServiceServer {
	wire.Build(
		loggerSet,
		service.NewUserService,
	)
	return pb.UserServiceServer(nil)
}

func InjectUserServiceMetrics() service.Metrics {
	wire.Build(
		service.MakeMetrics,
	)
	return service.Metrics{}
}
