package server

import (
	"context"
	"github.com/hashicorp/go-multierror"
	"github.com/sirupsen/logrus"
	"github.com/wuyuanyi135/mvprotos/mvcam"
	"github.com/wuyuanyi135/mvprotos/mvcamctrl"
	"github.com/wuyuanyi135/mvprotos/mvcgi"
	"google.golang.org/grpc"
)

type MicroVisionCGIServer struct {
	cameraCache     []*mvcam.DeviceInfo
	controllerCache []*mvcamctrl.SerialDeviceMapping

	cameraConnection *grpc.ClientConn
	cameraChannel    mvcam.MicroVisionCameraServiceClient

	controllerConnection *grpc.ClientConn
	controllerChannel    mvcamctrl.MicroVisionCameraControlServiceClient
}

func (s *MicroVisionCGIServer) ConnectDevices(ctx context.Context, req *mvcgi.ConnectDevicesRequest) (resp *mvcgi.ConnectDevicesResponse, err error) {
	err = multierror.Append(err, s.handleControllerConnection(ctx, req.ControllerConnection))
	err = multierror.Append(err, s.handleCameraConnection(ctx, req.CameraConnection))
	if err != nil {
		return
	}
	resp = &mvcgi.ConnectDevicesResponse{}
	return
}

func (s *MicroVisionCGIServer) ListDevices(ctx context.Context, req *mvcgi.ListDevicesRequest) (resp *mvcgi.ListDevicesResponse, err error) {
	// If the server urls are provided, update the connection to the designated server.
	if req.CameraUrl != "" {
		err = s.updateCameraConnection(req.CameraUrl)
		if err != nil {
			logrus.Errorf("Failed to connect to camera server %s: %s", req.CameraUrl, err.Error())
			return
		}
	}
	if req.ControllerUrl != "" {
		err = s.updateControllerConnection(req.ControllerUrl)
		if err != nil {
			logrus.Errorf("Failed to connect to controller server %s: %s", req.ControllerUrl, err.Error())
			return
		}
	}

	if !req.UseCache {
		// when cache update is requested, both caches will be updated even when another is giving error.
		localErr := s.updateCameraCache()
		if localErr != nil {
			err = multierror.Append(err, localErr)
		}
		localErr = s.updateControllerCache()
		if localErr != nil {
			err = multierror.Append(err, localErr)
		}
		if err != nil {
			return
		}
	}

	resp = &mvcgi.ListDevicesResponse{
		Devices:    s.cameraCache,
		Controller: s.controllerCache,
	}
	return
}

func (*MicroVisionCGIServer) GetVersion(ctx context.Context, req *mvcgi.GetVersionRequest) (resp *mvcgi.GetVersionResponse, err error) {
	resp = &mvcgi.GetVersionResponse{}
	err = nil
	return
}

func NewMicroVisionCGIServer() *MicroVisionCGIServer {
	return &MicroVisionCGIServer{}
}
