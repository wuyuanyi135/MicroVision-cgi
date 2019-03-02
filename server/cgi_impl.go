package server

import (
	"context"
	"github.com/hashicorp/go-multierror"
	"github.com/wuyuanyi135/mvprotos/mvcam"
	"github.com/wuyuanyi135/mvprotos/mvcamctrl"
	"github.com/wuyuanyi135/mvprotos/mvcgi"
	"golang.org/x/sync/semaphore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"sync"
	"time"
)

type MicroVisionCGIServer struct {
	BackendServerInterfaceSem *semaphore.Weighted
	CameraServerClient        *mvcam.MicroVisionCameraServiceClient
	CameraServerConnection    *grpc.ClientConn

	ControllerServerClient     *mvcamctrl.MicroVisionCameraControlServiceClient
	ControllerServerConnection *grpc.ClientConn
}

func NewMicroVisionCGIServer() *MicroVisionCGIServer {
	return &MicroVisionCGIServer{
		BackendServerInterfaceSem: semaphore.NewWeighted(1),
	}
}

func (s *MicroVisionCGIServer) GetVersion(ctx context.Context, req *mvcgi.GetVersionRequest) (resp *mvcgi.GetVersionResponse, err error) {
	resp = &mvcgi.GetVersionResponse{Version: "0.0.0"}
	return
}

// handle backend server connection control
func (s *MicroVisionCGIServer) BackendServerInterface(ctx context.Context, req *mvcgi.BackendServerInterfaceRequest) (resp *mvcgi.BackendServerInterfaceResponse, err error) {
	acquired := s.BackendServerInterfaceSem.TryAcquire(1)
	if !acquired {
		err = status.Error(codes.AlreadyExists, "a connection request is in progress. try again later")
		return
	}
	defer s.BackendServerInterfaceSem.Release(1)

	if req.CameraConnection != nil {
		errCameraConnection := s.handleCameraServerConnection(req.CameraConnection)
		if errCameraConnection != nil {
			err = multierror.Append(err, errCameraConnection)
		}
	}

	if req.ControllerConnection != nil {
		errControllerConnection := s.handleControllerServerConnection(req.ControllerConnection)
		if errControllerConnection != nil {
			err = multierror.Append(err, errControllerConnection)
		}
	}

	timeoutCtx1, _ := context.WithTimeout(ctx, 5*time.Second)
	timeoutCtx2, _ := context.WithTimeout(ctx, 5*time.Second)

	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		WaitUntilGrpcConnectionFinalState(timeoutCtx1, s.CameraServerConnection)
		wg.Done()
	}()

	go func() {
		WaitUntilGrpcConnectionFinalState(timeoutCtx2, s.ControllerServerConnection)
		wg.Done()
	}()
	wg.Wait()

	var (
		cameraServerStatus     mvcgi.ConnectionStatus
		controllerServerStatus mvcgi.ConnectionStatus
		localErr               error
	)

	if s.CameraServerConnection != nil {
		cameraServerStatus, localErr = TranslateGrpcStatusToMessage(s.CameraServerConnection.GetState(), "camera server")
		if localErr != nil {
			err = multierror.Append(localErr)
		}
	} else {
		cameraServerStatus = mvcgi.ConnectionStatus_DISCONNECTED
	}

	if s.ControllerServerConnection != nil {
		controllerServerStatus, localErr = TranslateGrpcStatusToMessage(s.ControllerServerConnection.GetState(), "controller server")
		if localErr != nil {
			err = multierror.Append(localErr)
		}
	} else {
		controllerServerStatus = mvcgi.ConnectionStatus_DISCONNECTED
	}


	// build the response
	resp = &mvcgi.BackendServerInterfaceResponse{
		CameraServer:     cameraServerStatus,
		ControllerServer: controllerServerStatus,
	}

	return
}

func (s *MicroVisionCGIServer) ListDevices(ctx context.Context, req *mvcgi.ListDevicesRequest) (resp *mvcgi.ListDevicesResponse, err error) {
	resp = &mvcgi.ListDevicesResponse{}
	if req.ListCamera != nil {
		var localErr error
		resp.Devices, localErr = s.handleListCamera(ctx, req.ListCamera)
		if localErr != nil {
			err = multierror.Append(err, localErr)
		}
	}

	if req.ListController != nil {
		var localErr error
		resp.Controller, localErr = s.handleListController(ctx, req.ListController)
		if localErr != nil {
			err = multierror.Append(err, localErr)
		}
	}
	return
}

func (s *MicroVisionCGIServer) DeviceInterface(ctx context.Context, req *mvcgi.DeviceInterfaceRequest) (resp *mvcgi.DeviceInterfaceResponse, err error) {
	resp = &mvcgi.DeviceInterfaceResponse{}
	if req.ControlCamera != nil {
		var localErr error
		resp.AffectedCamera, localErr = s.handleControlCamera(ctx, req.ControlCamera)
		if localErr != nil {
			err = multierror.Append(err, localErr)
		}
	}

	if req.ControlController != nil {
		var localErr error
		resp.AffectedController, localErr = s.handleControlController(ctx, req.ControlController)
		if localErr != nil {
			err = multierror.Append(err, localErr)
		}
	}
	return
}

func (s *MicroVisionCGIServer) DeviceParameterInterface(ctx context.Context, req *mvcgi.DeviceParameterInterfaceRequest) (resp *mvcgi.DeviceParameterInterfaceResponse, err error) {
	panic("implement me")
}

func (s *MicroVisionCGIServer) CameraStreaming(req *mvcgi.CameraStreamingRequest, stream mvcgi.MicroVisionCGI_CameraStreamingServer) (err error) {
	panic("implement me")
}

func (s *MicroVisionCGIServer) CameraCapturing(ctx context.Context, req *mvcgi.CameraCapturingRequest) (resp *mvcgi.CameraCapturingResponse, err error) {
	panic("implement me")
}


