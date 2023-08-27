package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/brickzzhang/grpc-helloworld/apigen/hello"
	"github.com/brickzzhang/grpc-helloworld/workshop/configger"
	"github.com/brickzzhang/grpc-helloworld/workshop/logger"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
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
	ctx = context.WithValue(ctx, "before_invoke", "bzz")
	header := metadata.New(map[string]string{"hello": "world"})
	header = metadata.Pairs("hello", "bzz")
	ctx = metadata.NewOutgoingContext(ctx, header)
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	// ctx, cancel := context.WithCancel(ctx)
	go func() {
		time.Sleep(1 * time.Second)
		cancel()
	}()
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

	// invoke with MaskField
	log.Printf("### ready to invoke Chatter2Hello ###\n")
	if err := client.FieldmaskTest(ctx); err != nil {
		logger.Error(ctx, "Chatter2Hello", zap.Error(err))
		return
	}
	log.Printf("\n\n")
}
