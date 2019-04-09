package devcon_test

import (
	"context"
	"fmt"
	"github.com/wuyuanyi135/MicroVisionCGI/server/devcon"
	"net"
	"testing"

	"github.com/wuyuanyi135/MicroVisionCGI/server"
	"github.com/wuyuanyi135/mvprotos/mvcam"
	"github.com/wuyuanyi135/mvprotos/mvcamctrl"
	"github.com/wuyuanyi135/mvprotos/mvcgi"
	"google.golang.org/grpc"
)

const Port = 30501

var client mvcgi.DeviceConnectionServiceClient

func TestMain(m *testing.M) {
	go StartTestServer()
	conn, err := grpc.Dial(fmt.Sprintf("0.0.0.0:%d", Port), grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	client = mvcgi.NewDeviceConnectionServiceClient(conn)

	m.Run()
}

func StartTestServer() {
	grpcServer := grpc.NewServer()
	cameraServerConn := server.BuildConnectionCameraServer()
	controllerServerConn := server.BuildConnectionControllerServer()

	cameraServer := mvcam.NewMicroVisionCameraServiceClient(cameraServerConn)
	controllerServer := mvcamctrl.NewMicroVisionCameraControlServiceClient(controllerServerConn)

	deviceConnectionService := devcon.NewDeviceConnectionServiceImpl(cameraServer, controllerServer)
	mvcgi.RegisterDeviceConnectionServiceServer(grpcServer, deviceConnectionService)
	address := fmt.Sprintf("0.0.0.0:%d", Port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		panic(err)
	}
	err = grpcServer.Serve(listener)
	if err != nil {
		panic(err)
	}

}

func TestDeviceDiscoveryServiceImpl_DiscoveryDevices(t *testing.T) {
	req := &mvcgi.DiscoveryDevicesRequest{DiscoverController: true, DiscoverCamera: true}
	response, err := client.DiscoveryDevices(context.Background(), req, grpc.WaitForReady(true))
	if err != nil {
		t.Error(err)
	}
	for k, v := range response.DiscoveredController {
		t.Logf("Found controller: #%d: id=%s; name=%s; connected=%t", k, v.Id, v.DisplayName, v.Connected)
	}
	for k, v := range response.DiscoveredCamera {
		t.Logf("Found camera: #%d: id=%s; name=%s; connected=%t", k, v.Id, v.DisplayName, v.Connected)
	}
}
