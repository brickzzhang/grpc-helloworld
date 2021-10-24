// Package handler gRPC gateway handler hook
package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// CustomizedResponse customized response
type CustomizedResponse struct {
	Code int `json:"code"`
	// Data returns response Data when status is ok
	Data interface{} `json:"data"`
	// Error returns Error message when status is not ok
	Error string `json:"error"`
}

const successFlag = "__success__"

// HTTPSuccHandler mux server http handler
func HTTPSuccHandler(ctx context.Context, w http.ResponseWriter, m proto.Message) error {
	mj := jsonpb.Marshaler{EmitDefaults: true}
	mjs, err := mj.MarshalToString(m)
	if err != nil {
		return err
	}
	resp := CustomizedResponse{
		Code:  200,
		Error: "",
	}
	if err = json.Unmarshal([]byte(mjs), &resp.Data); err != nil {
		return err
	}
	bs, err := json.Marshal(&resp)
	if err != nil {
		return err
	}
	return errors.New(successFlag + string(bs))
}

// HTTPErrorHandler mux server http handler
func HTTPErrorHandler(
	ctx context.Context, mux *runtime.ServeMux, m runtime.Marshaler, w http.ResponseWriter, r *http.Request, err error,
) {
	// success error
	raw := err.Error()
	if strings.HasPrefix(raw, successFlag) {
		raw = raw[len(successFlag):]
		_, _ = w.Write([]byte(raw))
		return
	}

	// normal error
	s, ok := status.FromError(err)
	if !ok {
		s = status.New(codes.Unknown, err.Error())
	}
	resp := CustomizedResponse{
		Code:  int(s.Code()),
		Data:  nil,
		Error: s.Message(),
	}
	bs, _ := json.Marshal(&resp)
	_, _ = w.Write(bs)
}
