package server

import (
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"github.com/wuyuanyi135/mvprotos/mvcgi"
	"google.golang.org/grpc"
)

func BuildServer() *grpcweb.WrappedGrpcServer {
	grpcServer := grpc.NewServer()
	cgiService := NewMicroVisionCGIServer()
	mvcgi.RegisterMicroVisionCGIServer(grpcServer, cgiService)
	wrappedGrpc := grpcweb.WrapServer(grpcServer)

	return wrappedGrpc
}
