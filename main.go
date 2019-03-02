package main

import (
	"github.com/gin-gonic/gin"
	"github.com/vulcand/oxy/forward"
	"github.com/vulcand/oxy/testutils"
	"github.com/wuyuanyi135/MicroVisionCGI/server"
	"github.com/wuyuanyi135/MicroVisionCGI/server/middleware"
	"google.golang.org/grpc/grpclog"
	"net/http"
)

func main() {
	wrappedGrpc := server.BuildServer()

	//// TODO: for dev. create static file serve for production.

	router := gin.Default()
	router.Use(middleware.GinGrpcWebMiddleware(wrappedGrpc))

	fwd, _ := forward.New()
	router.GET("/*proxy", func(context *gin.Context) {
		context.Request.URL = testutils.ParseURI("http://localhost:4200")
		fwd.ServeHTTP(context.Writer, context.Request)
	})

	if err := http.ListenAndServe(":8088", router); err != nil {
		grpclog.Fatalf("failed starting http2 server: %v", err)
	}
}