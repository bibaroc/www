package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/bibaroc/www/backend/deps"
	"github.com/bibaroc/www/backend/internal/user/pb"
	"github.com/bibaroc/www/backend/internal/user/service"
	"github.com/oklog/oklog/pkg/group"
	"google.golang.org/grpc"
)

func main() {
	logger := deps.InjectLogger()

	var (
		debugAddr = envStr("WWW_METRICS_ADDR", "0.0.0.0:8000")
		httpAddr  = envStr("WWW_HTTP_ADDR", "0.0.0.0:8001")
		grpcAddr  = envStr("WWW_GRPC_ADDR", "0.0.0.0:8002")
	)

	srv := deps.InjectUserService()
	if envBool("WWW_ENABLE_SERVICE_LOGGING", true) {
		srv = service.LoggingMiddleware(logger)(srv)
	}
	if envBool("WWW_ENABLE_SERVICE_METRICS", true) {
		srv = service.InstrumentingMiddleware(deps.InjectUserServiceMetrics())(srv)
	}

	var (
		endpoints   = service.NewEndpointSet(srv)
		httpHandler = service.NewHTTPHandler(endpoints, logger)
		grpcServer  = service.NewGRPCServer(endpoints, logger)
	)

	var g group.Group
	if envBool("WWW_ENABLE_SERVICE_METRICS", true) {
		debugListener, err := net.Listen("tcp", debugAddr)
		if err != nil {
			_ = logger.Log("transport", "debug/HTTP", "during", "Listen", "err", err)
			os.Exit(1)
		}
		g.Add(func() error {
			_ = logger.Log("transport", "debug/HTTP", "addr", debugAddr)
			return http.Serve(debugListener, http.DefaultServeMux)
		}, func(error) {
			debugListener.Close()
		})
	}
	if envBool("WWW_ENABLE_HTTP", true) {
		httpListener, err := net.Listen("tcp", httpAddr)
		if err != nil {
			_ = logger.Log("transport", "HTTP", "during", "Listen", "err", err)
			os.Exit(1)
		}
		g.Add(func() error {
			_ = logger.Log("transport", "HTTP", "addr", httpAddr)
			return http.Serve(httpListener, httpHandler)
		}, func(error) {
			httpListener.Close()
		})
	}
	if envBool("WWW_ENABLE_GRPC", true) {
		grpcListener, err := net.Listen("tcp", grpcAddr)
		if err != nil {
			_ = logger.Log("transport", "gRPC", "during", "Listen", "err", err)
			os.Exit(1)
		}
		g.Add(func() error {
			_ = logger.Log("transport", "gRPC", "addr", grpcAddr)
			baseServer := grpc.NewServer()
			pb.RegisterUserServiceServer(baseServer, grpcServer)
			return baseServer.Serve(grpcListener)
		}, func(error) {
			grpcListener.Close()
		})
	}
	{
		// This function just sits and waits for ctrl-C.
		cancelInterrupt := make(chan struct{})
		g.Add(func() error {
			c := make(chan os.Signal, 1)
			signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
			select {
			case sig := <-c:
				return fmt.Errorf("received signal %s", sig)
			case <-cancelInterrupt:
				return nil
			}
		}, func(error) {
			close(cancelInterrupt)
		})
	}
	_ = logger.Log("exit", g.Run())
}

func envStr(key, def string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return def
}
func envBool(key string, def bool) bool {
	if val, ok := os.LookupEnv(key); ok {
		b, err := strconv.ParseBool(val)
		if err != nil {
			return def
		}
		return b
	}
	return def
}
