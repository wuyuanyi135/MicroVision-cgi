package server

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/sirupsen/logrus"
	"github.com/wuyuanyi135/mvprotos/mvcam"
	"github.com/wuyuanyi135/mvprotos/mvcamctrl"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/status"
	"time"
)

// Update the connection /Reconnect to the given URL for the camera rpc server.
func (s *MicroVisionCGIServer) updateCameraConnection(serverUrl string) error {

	if s.cameraConnection != nil && serverUrl == s.cameraConnection.Target() {
		// no reconnection required
		logrus.Info("Reconnection to the camera is not required")
		return nil
	}

	conn, e := grpc.Dial(serverUrl, grpc.WithInsecure())
	if e != nil {
		return e
	}
	s.cameraConnection = conn

	client := mvcam.NewMicroVisionCameraServiceClient(conn)
	s.cameraChannel = client
	logrus.Infof("Reconnected to the camera: %s", serverUrl)

	return nil
}

// Update the connection /Reconnect to the given URL for the camera controller rpc server.
func (s *MicroVisionCGIServer) updateControllerConnection(serverUrl string) (err error) {
	if s.controllerConnection != nil && serverUrl == s.controllerConnection.Target() {
		// no reconnection required
		logrus.Info("Reconnection to the controller is not required")
		return nil
	}

	conn, e := grpc.Dial(serverUrl, grpc.WithInsecure())
	if e != nil {
		return e
	}
	s.controllerConnection = conn

	client := mvcamctrl.NewMicroVisionCameraControlServiceClient(conn)
	s.controllerChannel = client
	logrus.Infof("Reconnected to the controller: %s", serverUrl)

	return nil
}

// update the cached camera list
func (s *MicroVisionCGIServer) updateCameraCache() (err error) {
	err = s.cameraConnectionGuard()
	if err != nil {
		return err
	}

	state := s.cameraConnection.GetState()
	// wait for state change for the first connection
	for {
		if state == connectivity.Idle || state == connectivity.Connecting {
			ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
			s.cameraConnection.WaitForStateChange(ctx, state)
			state = s.cameraConnection.GetState()
		} else {
			break
		}
	}
	if state != connectivity.Ready {
		return status.Error(codes.Internal, fmt.Sprintf("Camera server not ready: %s", state))
	}

	response, err := s.cameraChannel.GetDevices(context.Background(), &mvcam.AdapterRequest{})
	if err != nil {
		return
	}
	s.cameraCache = response.Devices
	return
}

// update the cached controller list
func (s *MicroVisionCGIServer) updateControllerCache() (err error) {
	err = s.controllerConnectionGuard()
	if err != nil {
		return
	}

	state := s.controllerConnection.GetState()
	for {
		if state == connectivity.Idle || state == connectivity.Connecting {
			ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
			s.controllerConnection.WaitForStateChange(ctx, state)
			state = s.controllerConnection.GetState()
		} else {
			break
		}
	}

	if state != connectivity.Ready {
		return status.Error(codes.Internal, fmt.Sprintf("controller server not ready: %s", state))
	}

	response, err := s.controllerChannel.GetSerialDevices(context.Background(), &empty.Empty{})
	if err != nil {
		return
	}
	s.controllerCache = response.DeviceList
	return
}
