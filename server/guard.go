package server

import (
	"errors"
	"google.golang.org/grpc/connectivity"
)

func (s *MicroVisionCGIServer) guardCameraServerConnection() (err error) {
	if s.CameraServerClient == nil || s.CameraServerConnection == nil || s.CameraServerConnection.GetState() != connectivity.Ready {
		err = errors.New("Camera server is not connected!")
	}
	return
}
func (s *MicroVisionCGIServer) guardControllerServerConnection() (err error) {
	if s.ControllerServerClient == nil || s.ControllerServerConnection == nil || s.ControllerServerConnection.GetState() != connectivity.Ready {
		err = errors.New("Controller server is not connected!")
	}
	return
}