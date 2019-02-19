package middleware

import "github.com/improbable-eng/grpc-web/go/grpcweb"

type GrpcWebMiddleware struct {
	*grpcweb.WrappedGrpcServer
}

func NewGrpcWebMiddleware(wrappedGrpcServer *grpcweb.WrappedGrpcServer) *GrpcWebMiddleware {
	return &GrpcWebMiddleware{WrappedGrpcServer: wrappedGrpcServer}
}


