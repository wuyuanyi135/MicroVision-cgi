package server

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/wuyuanyi135/mvprotos/mvcgi"
)

func (s *MicroVisionCGIServer) handleControllerConnection(ctx context.Context, connection *mvcgi.ConnectDevicesRequest_ControllerConnection) (err error) {
	err = s.controllerConnectionGuard()
	if err != nil {
		return
	}

	switch connection.Action {
	case mvcgi.ConnectionAction_CONNECT:
		// connect to the controller
		_, err = s.controllerChannel.Connect(ctx, connection.ConnectionInfo)
		if err != nil {
			return
		}
		break
	case mvcgi.ConnectionAction_DISCONNECT:
		_, err = s.controllerChannel.Disconnect(ctx, &empty.Empty{})
		if err != nil {
			return
		}
		break
	default:
		break
	}
	return
}

func (s *MicroVisionCGIServer) handleCameraConnection(ctx context.Context, connection *mvcgi.ConnectDevicesRequest_CameraConnection) (err error) {
	err = s.cameraConnectionGuard()
	if err != nil {
		return
	}

	switch connection.Action {
	case mvcgi.ConnectionAction_NO_OP:
		return
	case mvcgi.ConnectionAction_CONNECT:
		// connect to the camera
		_, err = s.cameraChannel.OpenCamera(ctx, connection.ConnectionInfo)
		if err != nil {
			return
		}
		break
	case mvcgi.ConnectionAction_DISCONNECT:
		_, err = s.cameraChannel.ShutdownCamera(ctx, connection.ConnectionInfo)
		if err != nil {
			return
		}
	default:
		break
	}
	return
}
