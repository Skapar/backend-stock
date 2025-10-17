package server

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var notLogged = map[string]struct{}{}

func MiddlewareLoggingStream(logger *zap.Logger) grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		md, _ := metadata.FromIncomingContext(stream.Context())
		rid := md.Get("x-request-id")
		requestId := ""
		if len(rid) > 0 {
			requestId = rid[0]
		}

		// если метод есть в исключениях, то не логируем
		_, ok := notLogged[info.FullMethod]
		if !ok {
			l := logger.With(
				zap.String("method", info.FullMethod),
				zap.String("request_id", requestId),
			)
			l.Info("")
		}
		return handler(srv, stream)
	}
}

func MiddlewareLoggingUnary(logger *zap.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		md, _ := metadata.FromIncomingContext(ctx)
		rid := md.Get("x-request-id")
		requestId := ""
		if len(rid) > 0 {
			requestId = rid[0]
		}

		_, ok := notLogged[info.FullMethod]
		if !ok {
			l := logger.With(
				zap.String("method", info.FullMethod),
				zap.String("request_id", requestId),
				zap.Any("data", req),
			)
			l.Info("")
		}

		return handler(ctx, req)
	}
}
