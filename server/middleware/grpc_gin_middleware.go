package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"net/http"
)

func GinGrpcWebMiddleware(m *grpcweb.WrappedGrpcServer) gin.HandlerFunc {

	return func(context *gin.Context) {
		if m.IsAcceptableGrpcCorsRequest(context.Request) || m.IsGrpcWebRequest(context.Request) {
			context.Status(http.StatusOK) // prevent 404 in gin
			m.ServeHTTP(context.Writer, context.Request)
		} else {
			context.Next()
		}
	}
}
