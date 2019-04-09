package server

import (
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"github.com/wuyuanyi135/MicroVisionCGI/server/devicediscovery"
	"github.com/wuyuanyi135/MicroVisionCGI/server/devicepair"
	"github.com/wuyuanyi135/mvprotos/mvcam"
	"github.com/wuyuanyi135/mvprotos/mvcamctrl"
	"github.com/wuyuanyi135/mvprotos/mvcgi"
	"google.golang.org/grpc"
)

const CameraServerAddress = "localhost:"

func BuildServer() *grpcweb.WrappedGrpcServer {
	grpcServer := grpc.NewServer()
	// cgiService := NewMicroVisionCGIServer()
	// mvcgi.RegisterMicroVisionCGIServer(grpcServer, cgiService)

	cameraServerConn := BuildConnectionCameraServer()
	controllerServerConn := BuildConnectionControllerServer()

	cameraServer := mvcam.NewMicroVisionCameraServiceClient(cameraServerConn)
	controllerServer := mvcamctrl.NewMicroVisionCameraControlServiceClient(controllerServerConn)

	devicePairService := devicepair.NewDeviceServiceImpl()
	mvcgi.RegisterDevicePairServiceServer(grpcServer, devicePairService)
	deviceDiscoveryService := devicediscovery.NewDeviceDiscoveryServiceImpl(cameraServer, controllerServer)
	mvcgi.RegisterDeviceDiscoveryServiceServer(grpcServer, deviceDiscoveryService)
	wrappedGrpc := grpcweb.WrapServer(grpcServer)

	return wrappedGrpc
}
