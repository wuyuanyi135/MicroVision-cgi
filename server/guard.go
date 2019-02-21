package server

import "errors"

func (s *MicroVisionCGIServer) cameraConnectionGuard() error {
	if s.cameraConnection == nil || s.cameraChannel == nil {
		return errors.New("camera server not connected")
	}
	return nil
}

func (s *MicroVisionCGIServer) controllerConnectionGuard() error {
	if s.controllerConnection == nil || s.controllerChannel == nil {
		return errors.New("controller server not connected")
	}
	return nil
}
