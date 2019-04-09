package devcon

import (
	"context"
	"errors"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/wuyuanyi135/mvprotos/mvcam"
	"github.com/wuyuanyi135/mvprotos/mvcamctrl"
	"github.com/wuyuanyi135/mvprotos/mvcgi"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type DeviceConnectionServiceImpl struct {
	cameraServer     mvcam.MicroVisionCameraServiceClient
	controllerServer mvcamctrl.MicroVisionCameraControlServiceClient
}

func NewDeviceConnectionServiceImpl(
	cameraServer mvcam.MicroVisionCameraServiceClient,
	controllerServer mvcamctrl.MicroVisionCameraControlServiceClient,
) *DeviceConnectionServiceImpl {

	return &DeviceConnectionServiceImpl{
		cameraServer:     cameraServer,
		controllerServer: controllerServer,
	}
}

func (s *DeviceConnectionServiceImpl) Connect(ctx context.Context, req *mvcgi.ConnectionRequest) (resp *empty.Empty, err error) {
	switch req.Device.(type) {
	case *mvcgi.ConnectionRequest_CameraId:
		_, err = s.cameraServer.OpenCamera(ctx, &mvcam.IdRequest{Id: req.GetCameraId()}, grpc.WaitForReady(true))
		if err != nil {
			return
		}
		break
	case *mvcgi.ConnectionRequest_ControllerId:
		_, err = s.controllerServer.Connect(
			ctx,
			&mvcamctrl.ConnectRequest{
				DeviceIdentifier: &mvcamctrl.ConnectRequest_Path{
					Path: req.GetControllerId(),
				},
			},
			grpc.WaitForReady(true),
		)

		if err != nil {
			return
		}
		break
	default:
		err = errors.New("unknown device information provided")
		break
	}
	return
}

func (s *DeviceConnectionServiceImpl) Disconnect(ctx context.Context, req *mvcgi.ConnectionRequest) (resp *empty.Empty, err error) {
	switch req.Device.(type) {
	case *mvcgi.ConnectionRequest_CameraId:
		_, err = s.cameraServer.ShutdownCamera(ctx, &mvcam.IdRequest{Id: req.GetCameraId()}, grpc.WaitForReady(true))
		if err != nil {
			return
		}
		break
	case *mvcgi.ConnectionRequest_ControllerId:
		_, err = s.controllerServer.Disconnect(
			ctx,
			&mvcamctrl.ConnectRequest{
				DeviceIdentifier: &mvcamctrl.ConnectRequest_Path{
					Path: req.GetControllerId(),
				},
			},
			grpc.WaitForReady(true),
		)

		if err != nil {
			return
		}
		break
	default:
		err = errors.New("unknown device information provided")
		break
	}
	return
}

func (s *DeviceConnectionServiceImpl) DisconnectAll(ctx context.Context, req *mvcgi.DisconnectAllRequest) (resp *empty.Empty, err error) {
	return nil, status.Error(codes.Unimplemented, "DisconnectAll not implemented")
}
