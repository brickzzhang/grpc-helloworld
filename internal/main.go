package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	. "github.com/brickzzhang/grpc-helloworld/api"
	"github.com/brickzzhang/grpc-helloworld/internal/service"
	"github.com/brickzzhang/grpc-helloworld/workshop/configger"
	"github.com/brickzzhang/grpc-helloworld/workshop/logger"
)

// GRPCRegFunc Used while registering protoc generated gRPC registration function
type GRPCRegFunc func(server *grpc.Server)

func registerGRPCRegFunc(server *grpc.Server, funcs ...GRPCRegFunc) {
	for _, f := range funcs {
		f(server)
	}
}

func registerHelloWorldService(server *grpc.Server) {
	RegisterHelloWorldServiceServer(server, service.NewHelloWorldService())
}

// GatewayRegFunc Used while registering protoc generated gRPC gateway registration function
type GatewayRegFunc func(context.Context, *runtime.ServeMux) error

func registerGatewayRegFunc(ctx context.Context, mux *runtime.ServeMux, funcs ...GatewayRegFunc) error {
	var err error
	for _, g := range funcs {
		if err = g(ctx, mux); err != nil {
			return err
		}
	}

	return nil
}

func registerHelloWorldServiceGateway(ctx context.Context, mux *runtime.ServeMux) error {
	cfg := configger.ExtractConfiggerFromCtx(ctx)
	port := cfg.Get("grpc.grpcServerPort")

	opts := []grpc.DialOption{grpc.WithInsecure()}
	return RegisterHelloWorldServiceHandlerFromEndpoint(ctx, mux, fmt.Sprintf("localhost:%s", port), opts)
}

func startServer(ctx context.Context) (err error) {
	cfg := configger.ExtractConfiggerFromCtx(ctx)
	port := cfg.Get("grpc.grpcServerPort")

	// Create a listener on TCP port
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port.(string)))
	if err != nil {
		return fmt.Errorf("Failed to listen: %+v", err)
	}

	// Create a gRPC server object
	s := grpc.NewServer()
	// register grpc service to server
	registerGRPCRegFunc(s,
		registerHelloWorldService,
	)

	// Serve gRPC server
	logger.Info(ctx, fmt.Sprintf("Serving gRPC on localhost:%s", port.(string)))
	go func() {
		err = s.Serve(lis)
	}()

	return
}

func startGateway(ctx context.Context) (err error) {
	cfg := configger.ExtractConfiggerFromCtx(ctx)
	port := cfg.Get("grpc.grpcGatewayPort")

	mux := runtime.NewServeMux()
	err = registerGatewayRegFunc(ctx, mux,
		registerHelloWorldServiceGateway,
	)
	if err != nil {
		return fmt.Errorf("registerGateway func error: %v", err)
	}

	// Serve gRPC server
	logger.Info(ctx, fmt.Sprintf("GRPC gateway on localhost:%s", port.(string)))
	go func() {
		err = http.ListenAndServe(fmt.Sprintf(":%s", port.(string)), mux)
	}()
	return
}

func startProm(ctx context.Context) (err error) {
	cfg := configger.ExtractConfiggerFromCtx(ctx)
	port := cfg.Get("prometheus.prometheusPort")
	path := cfg.Get("prometheus.prometheusPath")

	router := mux.NewRouter()
	router.Path(path.(string)).Handler(promhttp.Handler())
	// Serve prometheus server
	logger.Info(ctx, fmt.Sprintf("Serving Prometheus requests on localhost:%s, route is %s", port.(string), path.(string)))
	go func() {
		err = http.ListenAndServe(fmt.Sprintf(":%s", port.(string)), router)
	}()

	return
}

func quitter(ctx context.Context) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		os.Interrupt,
	)

	sig := <-sigs

	logger.Info(ctx, fmt.Sprintf("%s signal received, server quit", sig.String()))
	os.Exit(0)
}

func main() {
	// initialize context
	ctx := context.Background()
	ctx = logger.NewLoggerToCtx(ctx, nil)
	ctx, err := configger.NewConfiggerToCtx(ctx)
	if err != nil {
		logger.Error(ctx, "config error", zap.Error(err))
		os.Exit(1)
	}

	// start server
	if err := startServer(ctx); err != nil {
		logger.Error(ctx, "start server error", zap.Error(err))
		os.Exit(1)
	}

	// start gateway
	if err := startGateway(ctx); err != nil {
		logger.Error(ctx, "start gateway error", zap.Error(err))
		os.Exit(1)
	}

	// start prometheus
	if err := startProm(ctx); err != nil {
		logger.Error(ctx, "start gateway error", zap.Error(err))
		os.Exit(1)
	}

	quitter(ctx)
}
