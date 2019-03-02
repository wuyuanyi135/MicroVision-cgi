package server

import (
	"context"
	"github.com/wuyuanyi135/mvprotos/mvcam"
	"github.com/wuyuanyi135/mvprotos/mvcamctrl"
	"github.com/wuyuanyi135/mvprotos/mvcgi"
)

func (s *MicroVisionCGIServer) handleControlCamera(ctx context.Context, options *mvcgi.DeviceInterfaceRequest_Options) (ret *mvcgi.DeviceInterfaceResponse_AffectedDevice, err error) {
	if options == nil {
		return
	}

	err = s.guardCameraServerConnection()
	if err != nil {
		return
	}

	switch options.Action {
	case mvcgi.ConnectionAction_CONNECT:
		_, err = (*s.CameraServerClient).OpenCamera(ctx, &mvcam.IdRequest{Id: options.Id})
		if err != nil {
			return
		}
		break
	case mvcgi.ConnectionAction_DISCONNECT:
		_, err = (*s.CameraServerClient).ShutdownCamera(ctx, &mvcam.IdRequest{Id: options.Id})
		if err != nil {
			return
		}
		break
	default:
		break
	}
	return
}

func (s *MicroVisionCGIServer) handleControlController(ctx context.Context, options *mvcgi.DeviceInterfaceRequest_Options) (ret *mvcgi.DeviceInterfaceResponse_AffectedDevice, err error) {
	if options == nil {
		return
	}

	err = s.guardControllerServerConnection()
	if err != nil {
		return
	}

	switch options.Action {
	case mvcgi.ConnectionAction_CONNECT:
		_, err = (*s.ControllerServerClient).Connect(ctx, &mvcamctrl.ConnectRequest{
			DeviceIdentifier: &mvcamctrl.ConnectRequest_Path{Path: options.Id},
			})
		if err != nil {
			return
		}
		break
	case mvcgi.ConnectionAction_DISCONNECT:
		_, err = (*s.ControllerServerClient).Disconnect(ctx, &mvcamctrl.ConnectRequest{
			DeviceIdentifier: &mvcamctrl.ConnectRequest_Path{Path: options.Id},
		})
		if err != nil {
			return
		}
		break
	default:
		break
	}
	return
}
