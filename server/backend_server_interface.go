package server

import (
	"context"
	"errors"
	"github.com/sirupsen/logrus"
	"github.com/wuyuanyi135/mvprotos/mvcam"
	"github.com/wuyuanyi135/mvprotos/mvcamctrl"
	"github.com/wuyuanyi135/mvprotos/mvcgi"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/status"
)

func (s *MicroVisionCGIServer) handleCameraServerConnection(req *mvcgi.BackendServerInterfaceRequest_Connection) (err error) {
	// reusable connect function
	connect := func() {
		conn, err := grpc.Dial(req.Url, grpc.WithInsecure())
		if err != nil {
			return
		}
		s.CameraServerConnection = conn
		client := mvcam.NewMicroVisionCameraServiceClient(conn)
		s.CameraServerClient = &client
		logrus.Infof("Connecting camera server to %s", req.Url)
	}
	disconnect := func() {
		defer func() {
			s.CameraServerClient = nil
			logrus.Infof("Disconnected camera server")
		}()
		if s.CameraServerConnection != nil {
			_ = s.CameraServerConnection.Close()
		}

	}

	// wants to perform action on the camera server
	switch req.Action {
	case mvcgi.ConnectionAction_CONNECT:
		// is there connection information?
		if req.Url == "" {
			err = status.Error(codes.InvalidArgument, "Camera server url not provided")
			return
		}

		// is it connected?
		if s.CameraServerConnection != nil && s.CameraServerConnection.GetState() == connectivity.Ready {
			// looks connected. Only accept connecting to other server
			if req.Url != s.CameraServerConnection.Target() {
				connect()
			}
		} else {
			// not connected or not ready or error, connect to the specified url anyway
			connect()
		}
		break
	case mvcgi.ConnectionAction_DISCONNECT:
		disconnect()
		break
	default:
		break
	}
	return
}

func (s *MicroVisionCGIServer) handleControllerServerConnection(req *mvcgi.BackendServerInterfaceRequest_Connection) (err error) {
	// reusable connect function
	connect := func() {
		conn, err := grpc.Dial(req.Url, grpc.WithInsecure())
		if err != nil {
			return
		}
		s.ControllerServerConnection = conn
		client := mvcamctrl.NewMicroVisionCameraControlServiceClient(conn)
		s.ControllerServerClient = &client
		logrus.Infof("Connecting controller server to %s", req.Url)

	}
	disconnect := func() {
		defer func() {
			s.ControllerServerClient = nil
			logrus.Info("Disconnecting controller server")
		}()
		if s.ControllerServerConnection != nil {
			_ = s.ControllerServerConnection.Close()
		}

	}

	// wants to perform action on the controller server
	switch req.Action {
	case mvcgi.ConnectionAction_CONNECT:
		// is there connection information?
		if req.Url == "" {
			err = status.Error(codes.InvalidArgument, "Controller server url not provided")
			return
		}

		// is it connected?
		if s.ControllerServerConnection != nil && s.ControllerServerConnection.GetState() == connectivity.Ready {
			// looks connected. Only accept connecting to other server
			if req.Url != s.ControllerServerConnection.Target() {
				connect()
			}
		} else {
			// not connected or not ready or error, connect to the specified url anyway
			connect()
		}
		break
	case mvcgi.ConnectionAction_DISCONNECT:
		disconnect()
		break
	default:
		break
	}
	return
}

func WaitUntilGrpcConnectionFinalState(ctx context.Context, conn *grpc.ClientConn) (state connectivity.State) {
	if conn == nil {
		return connectivity.Shutdown
	}
	for {
		state = conn.GetState()
		switch state {
		case connectivity.Connecting:
			fallthrough
		case connectivity.Idle:
			conn.WaitForStateChange(ctx, state)
			break
		case connectivity.TransientFailure:
			fallthrough
		case connectivity.Ready:
			fallthrough
		case connectivity.Shutdown:
			return
		}
	}
}

func WaitAllContext(ctx []context.Context) {
	for _, c := range ctx {
		<-c.Done()
	}
}

func TranslateGrpcStatusToMessage(status connectivity.State, server string) (ret mvcgi.ConnectionStatus, err error){
	switch status {
	case connectivity.TransientFailure:
		err = errors.New("Connection failed to " + server)
		return
	case connectivity.Shutdown:
		ret = mvcgi.ConnectionStatus_DISCONNECTED
		return
	case connectivity.Ready:
		ret = mvcgi.ConnectionStatus_CONNECTED
		return
	default:
		err = errors.New("Unexpected status: " + status.String())
		return
	}
}