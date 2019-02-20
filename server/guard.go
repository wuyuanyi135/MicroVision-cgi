package server

import "errors"

func (s *MicroVisionCGIServer) cameraConnectionGuard() error {
	if s.cameraConnection != nil && s.cameraChannel != nil {
		return errors.New("camera server not connected")
	}
	return nil
}

func (s *MicroVisionCGIServer) controllerConnectionGuard() error {
	if s.controllerConnection != nil && s.controllerChannel != nil {
		return errors.New("controller server not connected")
	}
	return nil
}

func (s *MicroVisionCGIServer) cameraOpenGuard() error {
	if !s.cameraOpened {
		return errors.New("camera is not opened")
	}
	return nil
}

func (s *MicroVisionCGIServer) controllerOpenGuard() error {
	if !s.controllerOpened {
		return errors.New("controller is not opened")
	}
	return nil
}
