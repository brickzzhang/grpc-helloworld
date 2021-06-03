package service

import (
	"context"

	. "github.com/brickzzhang/grpc-helloworld/api"
)

type HelloWorldService struct {
	UnimplementedHelloWorldServiceServer
}

func NewHelloWorldService() *HelloWorldService {
	service := &HelloWorldService{}
	return service
}

func (service *HelloWorldService) SayHello(ctx context.Context, request *SayHelloRequest) (*SayHelloResponse, error) {
	return &SayHelloResponse{Message: request.Message}, nil
}
