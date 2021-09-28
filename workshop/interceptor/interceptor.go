// Package interceptor interceptor demo
package interceptor

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// KeyType key type
type KeyType string

// OldKey old key injected into context by requester
var OldKey KeyType = "old-key"

// NewKey copy of OldKey injected into context by interceptor
var NewKey KeyType = "new-key"

// ExtractMetadata extract NewKey from context and return
func ExtractMetadata(ctx context.Context, key KeyType) (interface{}, bool) {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		v := md.Get(string(key))
		if len(v) > 0 {
			return v[0], true
		}

		return nil, true
	}
	return nil, false
}

// UnaryServerInterceptor customized unary server interceptor
// duplicates a new header field using old one, add it to context and return
func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			var value []string
			if v, ok := md[string(OldKey)]; ok {
				value = v
			}
			md1 := metadata.Pairs(string(NewKey), value[0])
			md2 := metadata.Join(md, md1)
			ctx = metadata.NewIncomingContext(ctx, md2)
		}

		return handler(ctx, req)
	}
}
