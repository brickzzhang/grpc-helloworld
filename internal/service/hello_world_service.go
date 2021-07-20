// Package service provide helloworld service
package service

import (
	"context"

	"github.com/brickzzhang/grpc-helloworld/apigen/hello"
)

// HelloWorldService helloworld service
type HelloWorldService struct {
	hello.UnimplementedHelloWorldServiceServer
}

// NewHelloWorldService new helloworld service
func NewHelloWorldService() *HelloWorldService {
	service := &HelloWorldService{}
	return service
}

// SayHello echo nil
func (service *HelloWorldService) SayHello(
	ctx context.Context, request *hello.SayHelloRequest) (*hello.SayHelloResponse, error) {
	return &hello.SayHelloResponse{Message: request.Message}, nil
}
