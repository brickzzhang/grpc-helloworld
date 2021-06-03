// Package service provide helloworld service
package service

import (
	"context"

	api "github.com/brickzzhang/grpc-helloworld/api"
)

// HelloWorldService helloworld service
type HelloWorldService struct {
	api.UnimplementedHelloWorldServiceServer
}

// NewHelloWorldService new helloworld service
func NewHelloWorldService() *HelloWorldService {
	service := &HelloWorldService{}
	return service
}

// SayHello echo nil
func (service *HelloWorldService) SayHello(
	ctx context.Context, request *api.SayHelloRequest) (*api.SayHelloResponse, error) {
	return &api.SayHelloResponse{Message: request.Message}, nil
}
