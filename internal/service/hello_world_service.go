// Package service provide helloworld service
package service

import (
	"context"
	"fmt"
	"io"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/brickzzhang/grpc-helloworld/apigen/hello"
	"github.com/brickzzhang/grpc-helloworld/workshop/interceptor"
	fieldmaskutils "github.com/mennanov/fieldmask-utils"
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
func (s *HelloWorldService) SayHello(
	ctx context.Context, r *hello.SayHelloRequest,
) (*hello.SayHelloResponse, error) {
	if v, ok := interceptor.ExtractMetadata(ctx, interceptor.NewKey); ok {
		log.Printf("value injected into the header: %+v", v)
	}
	return &hello.SayHelloResponse{Message: r.Message}, nil
}

// HelloChatter say hello continuously on server side
func (s *HelloWorldService) HelloChatter(
	r *hello.HelloChatterRequest, stream hello.HelloWorldService_HelloChatterServer,
) error {
	for i := 0; i < 60; i++ {
		if err := stream.Send(
			&hello.HelloChatterResponse{
				Seq:     int32(i),
				Message: fmt.Sprintf("echo %s %ded time", r.Message, i),
			},
		); err != nil {
			return err
		}

		time.Sleep(1 * time.Second)
	}

	return nil
}

// Chatter2Hello say hello continuously on client side
func (s *HelloWorldService) Chatter2Hello(
	stream hello.HelloWorldService_Chatter2HelloServer,
) error {
	var totalSeq int32

	for {
		r, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&hello.Chatter2HelloResponse{
				Total: totalSeq,
			})
		}
		if err != nil {
			return err
		}

		log.Printf("%ded message received: %s\n", r.Seq, r.Message)

		totalSeq++
	}
}

// Chatter2Chatter chatter to chatter talk
func (s *HelloWorldService) Chatter2Chatter(
	stream hello.HelloWorldService_Chatter2ChatterServer,
) error {
	wg := &sync.WaitGroup{}
	wg.Add(2)

	// send goroutine
	go func() {
		defer wg.Done()

		total := 30
		for i := 0; i < total; i++ {
			if err := stream.Send(&hello.Chatter2ChatterResponse{
				Seq:     int32(i),
				Message: fmt.Sprintf("seq: %d, %d in total, hello client", i, total),
			}); err != nil {
				log.Printf("Chatter2Chatter send message error: %+v", err)
			}
			time.Sleep(1 * time.Second)
		}
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
				log.Printf("Chatter2Chatter receive message error: %+v", err)
				return
			}
			log.Printf("%ded message received: %s\n", in.Seq, in.Message)
		}
	}()

	wg.Wait()
	return nil
}

// GetTest get test for body.
func (s *HelloWorldService) GetTest(
	ctx context.Context, req *hello.GetTestReq,
) (*hello.GetTestRes, error) {
	return &hello.GetTestRes{
		Msg: req.Msg,
	}, nil
}

// FieldmaskTest fieldmask with part body.
func (s *HelloWorldService) FieldmaskTest(
	ctx context.Context, req *hello.FieldmaskTestReq,
) (*hello.FieldmaskTestRes, error) {
	full := &hello.FieldmaskTestRes{
		Msg: "reply",
		Nested: &hello.Nested{
			Attr1: "attr1",
			Attr2: "attr2",
		},
	}
	res := &hello.FieldmaskTestRes{Nested: &hello.Nested{}}
	log.Printf("!!! mask: %+v, req: %+v", req.FieldMask, req)
	if req.FieldMask != nil {
		log.Printf("!!! paths: %+v", req.FieldMask.Paths)
		mask, err := fieldmaskutils.MaskFromPaths(req.FieldMask.Paths, func(s string) string {
			s = strings.Replace(s, "_", " ", -1)
			s = strings.Title(s)
			return strings.Replace(s, " ", "", -1)
		})
		if err != nil {
			log.Printf("mask error: %+v", err)
			return nil, err
		}
		fieldmaskutils.StructToStruct(mask, full.Nested, res.Nested)
	}

	return res, nil
}

// FieldmaskTestWithBodyStar fieldmask with body *.
func (s *HelloWorldService) FieldmaskTestWithBodyStar(
	ctx context.Context, req *hello.FieldmaskTestReq,
) (*hello.FieldmaskTestRes, error) {
	full := &hello.FieldmaskTestRes{
		Msg: "reply",
		Nested: &hello.Nested{
			Attr1: "attr1",
			Attr2: "attr2",
		},
	}
	res := &hello.FieldmaskTestRes{Nested: &hello.Nested{}}
	log.Printf("!!! mask: %+v, req: %+v", req.FieldMask, req)
	if req.FieldMask != nil {
		log.Printf("!!! paths: %+v", req.FieldMask.Paths)
		mask, err := fieldmaskutils.MaskFromPaths(req.FieldMask.Paths, func(s string) string {
			s = strings.Replace(s, "_", " ", -1)
			s = strings.Title(s)
			return strings.Replace(s, " ", "", -1)
		})
		if err != nil {
			log.Printf("mask error: %+v", err)
			return nil, err
		}
		fieldmaskutils.StructToStruct(mask, full, res)
	}

	return res, nil
}
