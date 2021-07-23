package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/brickzzhang/grpc-helloworld/apigen/hello"
	"github.com/brickzzhang/grpc-helloworld/workshop/configger"
	"github.com/brickzzhang/grpc-helloworld/workshop/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	// initialize context
	ctx := context.Background()
	ctx = logger.NewLoggerToCtx(ctx, nil)
	ctx, err := configger.NewConfiggerToCtx(ctx)
	if err != nil {
		logger.Error(ctx, "config error", zap.Error(err))
		os.Exit(1)
	}

	// establish connection
	cfg := configger.ExtractConfiggerFromCtx(ctx)
	port := cfg.Get("grpc.grpcServerPort")
	conn, err := grpc.Dial(fmt.Sprintf(":%s", port), grpc.WithInsecure())
	if err != nil {
		logger.Error(ctx, "grpc client dial error", zap.Error(err))
		return
	}
	defer func() {
		_ = conn.Close()
	}()

	// init client
	client := HelloWorldClient{
		Client: hello.NewHelloWorldServiceClient(conn),
	}

	// invoke SayHello
	log.Printf("### ready to invoke SayHello ###\n")
	if err := client.SayHello(ctx, &hello.SayHelloRequest{
		Message: "hello world",
	}); err != nil {
		logger.Error(ctx, "SayHello", zap.Error(err))
		return
	}
	log.Printf("\n\n")

	// invoke HelloChatter
	log.Printf("### ready to invoke HelloChatter ###\n")
	if err := client.HelloChatter(ctx, &hello.HelloChatterRequest{
		Message: "reply stream to me",
	}); err != nil {
		logger.Error(ctx, "HelloChatter", zap.Error(err))
		return
	}
	log.Printf("\n\n")

	// invoke Chatter2Hello
	log.Printf("### ready to invoke Chatter2Hello ###\n")
	if err := client.Chatter2Hello(ctx); err != nil {
		logger.Error(ctx, "Chatter2Hello", zap.Error(err))
		return
	}
	log.Printf("\n\n")

	// invoke Chatter2Chatter
	log.Printf("### ready to invoke Chatter2Chatter ###\n")
	if err := client.Chatter2Chatter(ctx); err != nil {
		logger.Error(ctx, "Chatter2Chatter", zap.Error(err))
		return
	}
	log.Printf("\n\n")
}
