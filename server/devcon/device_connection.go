package devcon

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/hashicorp/go-multierror"
	"github.com/wuyuanyi135/mvprotos/mvcam"
	"github.com/wuyuanyi135/mvprotos/mvcamctrl"
	"github.com/wuyuanyi135/mvprotos/mvcgi"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"sync"
)

type DeviceConnectionServiceImpl struct {
	cameraServer     mvcam.MicroVisionCameraServiceClient
	controllerServer mvcamctrl.MicroVisionCameraControlServiceClient
	cache            mvcgi.DiscoveryDevicesResponse
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
	resp = &empty.Empty{}
	wg := sync.WaitGroup{}
	e := &multierror.Error{}
	if req.CameraId != "" {
		wg.Add(1)
		go func() {
			defer wg.Done()

			_, err := s.cameraServer.OpenCamera(ctx, &mvcam.IdRequest{Id: req.CameraId}, grpc.WaitForReady(true))
			if err != nil {
				e = multierror.Append(e, err)
				return
			}
			// update cache
			for k, v := range s.cache.DiscoveredCamera {
				if v.Id == req.CameraId {
					s.cache.DiscoveredCamera[k].Connected = true
				}
			}
		}()
	}
	if req.ControllerId != "" {
		wg.Add(1)
		go func() {
			defer wg.Done()

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
				e = multierror.Append(e, err)
				return
			}

			// update cache
			for k, v := range s.cache.DiscoveredController {
				if v.Id == req.ControllerId {
					s.cache.DiscoveredController[k].Connected = true
				}
			}
		}()
	}

	wg.Wait()
	if e.ErrorOrNil() != nil {
		err = e
	}
	return
}

func (s *DeviceConnectionServiceImpl) Disconnect(ctx context.Context, req *mvcgi.ConnectionRequest) (resp *empty.Empty, err error) {
	resp = &empty.Empty{}
	wg := sync.WaitGroup{}
	e := &multierror.Error{}
	if req.CameraId != "" {
		wg.Add(1)
		go func() {
			defer wg.Done()

			_, err := s.cameraServer.ShutdownCamera(ctx, &mvcam.IdRequest{Id: req.CameraId}, grpc.WaitForReady(true))
			if err != nil {
				e = multierror.Append(e, err)
			}
			// update cache
			for k, v := range s.cache.DiscoveredCamera {
				if v.Id == req.CameraId {
					s.cache.DiscoveredCamera[k].Connected = false
				}
			}
		}()
	}
	if req.ControllerId != "" {
		wg.Add(1)
		go func() {
			defer wg.Done()

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
				e = multierror.Append(e, err)
			}
			// update cache
			for k, v := range s.cache.DiscoveredController {
				if v.Id == req.ControllerId {
					s.cache.DiscoveredController[k].Connected = false
				}
			}
		}()
	}

	wg.Wait()
	if e.ErrorOrNil() != nil {
		err = e
	}
	return
}

func (s *DeviceConnectionServiceImpl) DisconnectAll(ctx context.Context, req *mvcgi.DisconnectAllRequest) (resp *empty.Empty, err error) {
	resp = &empty.Empty{}
	return resp, status.Error(codes.Unimplemented, "DisconnectAll not implemented")
}
