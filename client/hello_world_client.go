package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"sync"
	"time"

	"github.com/brickzzhang/grpc-helloworld/apigen/hello"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

// HelloWorldClient helloworld client
type HelloWorldClient struct {
	Client hello.HelloWorldServiceClient
}

// SayHello invoke grpc SayHello
func (c *HelloWorldClient) SayHello(ctx context.Context, req *hello.SayHelloRequest) error {
	cc := c.Client
	res, err := cc.SayHello(ctx, req)
	if err != nil {
		return err
	}

	log.Printf("response is: %+v", res)
	return nil
}

// HelloChatter invoke grpc HelloChatter
func (c *HelloWorldClient) HelloChatter(ctx context.Context, req *hello.HelloChatterRequest) error {
	cc := c.Client
	stream, err := cc.HelloChatter(ctx, req)
	if err != nil {
		return err
	}
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		log.Printf("seq: %d, msg: %s", res.Seq, res.Message)
	}
}

// Chatter2Hello invoke grpc Chatter2Hello
func (c *HelloWorldClient) Chatter2Hello(ctx context.Context) error {
	cc := c.Client
	stream, err := cc.Chatter2Hello(ctx)
	if err != nil {
		return err
	}

	for i := 0; i < 60; i++ {
		if err = stream.Send(&hello.Chatter2HelloRequest{
			Seq:     int32(i),
			Message: fmt.Sprintf("seq: %d, hello server", i),
		}); err != nil {
			return err
		}

		time.Sleep(1 * time.Second)
	}

	reply, err := stream.CloseAndRecv()
	if err != nil {
		return err
	}
	log.Printf("reply: %+v", reply)
	return nil
}

// Chatter2Chatter invoke grpc Chatter2Chatter
func (c *HelloWorldClient) Chatter2Chatter(ctx context.Context) error {
	cc := c.Client
	stream, err := cc.Chatter2Chatter(ctx)
	if err != nil {
		return err
	}

	wg := &sync.WaitGroup{}
	wg.Add(2)

	// send goroutine
	go func() {
		defer wg.Done()

		total := 60
		for i := 0; i < total; i++ {
			if err := stream.Send(&hello.Chatter2ChatterRequest{
				Seq:     int32(i),
				Message: fmt.Sprintf("seq: %d, %d in total, hello server", i, total),
			}); err != nil {
				log.Printf("Chatter2Chatter send message error: %+v", err)
			}
			time.Sleep(1 * time.Second)
		}
		_ = stream.CloseSend()
	}()

	// receive goroutine
	go func() {
		defer wg.Done()

		for {
			in, err := stream.Recv()
			if err == io.EOF {
				return
			}
			if err != nil {
				log.Printf("Chatter2Chatter receive error: %+v", err)
				return
			}
			log.Printf("%ded message received: %s\n", in.Seq, in.Message)
		}

	}()

	wg.Wait()
	return nil
}

// FieldmaskTest api
func (c *HelloWorldClient) FieldmaskTest(ctx context.Context) error {
	cc := c.Client
	req := &hello.FieldmaskTestReq{
		Msg: "test msg",
		Nested: &hello.Nested{
			Attr1: "attr1",
		},
	}
	fm, err := fieldmaskpb.New(req, "msg")
	if err != nil {
		log.Printf("construct field_mask error: %+v", err)
		return err
	}
	req.FieldMask = fm
	res, err := cc.FieldmaskTest(ctx, req)
	if err != nil {
		log.Printf("invoke error: %+v", err)
		return err
	}

	log.Printf("res: %+v", res)
	return nil
}
