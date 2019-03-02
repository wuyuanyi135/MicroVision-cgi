package server_test

import (
	"context"
	"fmt"
	"github.com/wuyuanyi135/MicroVisionCGI/server"
	"github.com/wuyuanyi135/mvprotos/mvcgi"
	"google.golang.org/grpc"
	"net"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	grpcServer := grpc.NewServer()
	cgiService := server.NewMicroVisionCGIServer()
	mvcgi.RegisterMicroVisionCGIServer(grpcServer, cgiService)
	var listener net.Listener
	var err error
	go func() {
		listener, err = net.Listen("tcp", ":9901")
		if err != nil {
			panic(err)
		}
		err = grpcServer.Serve(listener)
		if err != nil {
			fmt.Println(err)
		}

	}()
	retCode := m.Run()
	err = listener.Close()
	if err != nil {
		panic(err)
	}
	os.Exit(retCode)
}

func GetClient() mvcgi.MicroVisionCGIClient {
	clientConn, err := grpc.Dial("localhost:9901", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	client := mvcgi.NewMicroVisionCGIClient(clientConn)
	return client
}

func OpenCtrollerServer(client mvcgi.MicroVisionCGIClient)  {
	_, err := client.BackendServerInterface(context.Background(), &mvcgi.BackendServerInterfaceRequest{
		ControllerConnection: &mvcgi.BackendServerInterfaceRequest_Connection{
			Action: mvcgi.ConnectionAction_CONNECT,
			Url:    "localhost:3050",
		},
	})
	if err != nil {
		panic(err)
	}
}
func OpenCameraServer(client mvcgi.MicroVisionCGIClient)  {
	_, err := client.BackendServerInterface(context.Background(), &mvcgi.BackendServerInterfaceRequest{
		CameraConnection: &mvcgi.BackendServerInterfaceRequest_Connection{
			Action: mvcgi.ConnectionAction_CONNECT,
			Url:    "localhost:5074",
		},
	})
	if err != nil {
		panic(err)
	}
}
func CloseServers(client mvcgi.MicroVisionCGIClient) {
	_, err := client.BackendServerInterface(context.Background(), &mvcgi.BackendServerInterfaceRequest{
		ControllerConnection: &mvcgi.BackendServerInterfaceRequest_Connection{
			Action: mvcgi.ConnectionAction_DISCONNECT,
		},
		CameraConnection: &mvcgi.BackendServerInterfaceRequest_Connection{
			Action: mvcgi.ConnectionAction_DISCONNECT,
		},
	})
	if err != nil {
		panic(err)
	}
}