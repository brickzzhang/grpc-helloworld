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
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	hello "github.com/brickzzhang/grpc-helloworld/apigen/hello"
	"github.com/brickzzhang/grpc-helloworld/internal/service"
	"github.com/brickzzhang/grpc-helloworld/workshop/configger"
	"github.com/brickzzhang/grpc-helloworld/workshop/handler"
	"github.com/brickzzhang/grpc-helloworld/workshop/interceptor"
	"github.com/brickzzhang/grpc-helloworld/workshop/logger"
	"github.com/brickzzhang/grpc-helloworld/workshop/swagger"
)

// GRPCRegFunc Used while registering protoc generated gRPC registration function
type GRPCRegFunc func(server *grpc.Server)

func registerGRPCRegFunc(server *grpc.Server, funcs ...GRPCRegFunc) {
	for _, f := range funcs {
		f(server)
	}
}

func registerHelloWorldService(server *grpc.Server) {
	hello.RegisterHelloWorldServiceServer(server, service.NewHelloWorldService())
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
	return hello.RegisterHelloWorldServiceHandlerFromEndpoint(ctx, mux, fmt.Sprintf("localhost:%s", port), opts)
}

func startServer(ctx context.Context) (err error) {
	cfg := configger.ExtractConfiggerFromCtx(ctx)
	port := cfg.Get("grpc.grpcServerPort")

	// Create a listener on TCP port
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port.(string)))
	if err != nil {
		return fmt.Errorf("failed to listen: %+v", err)
	}

	// log grpc library internals
	// grpc_zap.ReplaceGrpcLogger(logger.GetDefaultZapLogger())

	// log payloads for all requests and responses
	alwaysLoggingDeciderServer := func(ctx context.Context, fullMethodName string, servingObject interface{}) bool {
		return true
	}

	// Create a gRPC server object
	s := grpc.NewServer(
		grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(
				grpc_recovery.UnaryServerInterceptor(),
				grpc_ctxtags.UnaryServerInterceptor(
					grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor),
				),
				grpc_zap.UnaryServerInterceptor(logger.GetDefaultZapLogger()),
				grpc_zap.PayloadUnaryServerInterceptor(
					logger.GetDefaultZapLogger(), alwaysLoggingDeciderServer,
				),
				grpc_opentracing.UnaryServerInterceptor(),
				grpc_prometheus.UnaryServerInterceptor,
				interceptor.UnaryServerInterceptor(),
			),
		),
		grpc.StreamInterceptor(
			grpc_middleware.ChainStreamServer(
				grpc_recovery.StreamServerInterceptor(),
				grpc_ctxtags.StreamServerInterceptor(
					grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor),
				),
				grpc_zap.StreamServerInterceptor(logger.GetDefaultZapLogger()),
				grpc_zap.PayloadStreamServerInterceptor(
					logger.GetDefaultZapLogger(), alwaysLoggingDeciderServer,
				),
				grpc_opentracing.StreamServerInterceptor(),
				grpc_prometheus.StreamServerInterceptor,
			),
		),
	)
	// register grpc service to server
	registerGRPCRegFunc(s,
		registerHelloWorldService,
	)

	// todo: debug mode.
	reflection.Register(s)

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

	mux := runtime.NewServeMux(
		runtime.WithForwardResponseOption(handler.HTTPSuccHandler),
		runtime.WithProtoErrorHandler(handler.HTTPErrorHandler),
	)
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

func startSwagger(ctx context.Context) (err error) {
	cfg := configger.ExtractConfiggerFromCtx(ctx)
	port := cfg.Get("swagger.swaggerWebPort")
	grpcGwPort := cfg.Get("grpc.grpcGatewayPort")
	// assign the web route
	value, ok := cfg.Get("swagger.swaggerWebPath").(string)
	if !ok {
		return fmt.Errorf("assert error: %+v", value)
	}
	swagger.SwWebRoute = value
	enabled, ok := cfg.Get("swagger.swaggerEnabled").(bool)
	if !ok {
		return fmt.Errorf("assert error: %+v", enabled)
	}
	if !enabled {
		return nil
	}

	mux := http.NewServeMux()
	// http://localhost:8082/sw handler
	mux.Handle(swagger.SwWebRoute, http.StripPrefix(swagger.SwWebRoute, swagger.WebHandler()))
	// http://localhost:8082/swagger/swagger/application/v1/application_service.swagger.json handler
	mux.HandleFunc(swagger.SwJSONRoute, swagger.WebJSONHandler)
	// http://localhost:8092/v1/helloworld handler
	mux.Handle(swagger.SwaggerGatewayRoute, swagger.Forward2GprcGatewayHandler(grpcGwPort.(string)))

	// Serve swagger server
	logger.Info(ctx, fmt.Sprintf("swagger on localhost:%s%s", port.(string), swagger.SwWebRoute))
	go func() {
		err = http.ListenAndServe(fmt.Sprintf(":%s", port.(string)), mux)
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

	// start swagger, error doesn't matter
	if err := startSwagger(ctx); err != nil {
		logger.Error(ctx, "start swagger error", zap.Error(err))
	}

	// start prometheus
	if err := startProm(ctx); err != nil {
		logger.Error(ctx, "start gateway error", zap.Error(err))
		os.Exit(1)
	}

	quitter(ctx)
}
