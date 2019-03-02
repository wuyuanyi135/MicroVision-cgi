package server

import (
	"context"
	"github.com/wuyuanyi135/mvprotos/mvcam"
	"github.com/wuyuanyi135/mvprotos/mvcamctrl"
	"github.com/wuyuanyi135/mvprotos/mvcgi"
)

func (s *MicroVisionCGIServer) handleListCamera(ctx context.Context, options *mvcgi.ListDevicesRequest_Options) (ret []*mvcam.DeviceInfo, err error) {
	err = s.guardCameraServerConnection()
	if err != nil {
		return
	}

	response, err := (*s.CameraServerClient).GetDevices(ctx, &mvcam.GetDevicesRequest{UseCache:options != nil && options.UseCache})
	if err != nil {
		return
	}
	ret = response.Devices
	return
}
func (s *MicroVisionCGIServer) handleListController(ctx context.Context, options *mvcgi.ListDevicesRequest_Options) (ret []*mvcamctrl.SerialDeviceMapping, err error) {
	err = s.guardControllerServerConnection()
	if err != nil {
		return
	}

	response, err := (*s.ControllerServerClient).GetSerialDevices(ctx, &mvcamctrl.GetSerialDevicesRequest{UseCache: options != nil && options.UseCache})
	if err != nil {
		return
	}
	ret = response.DeviceList
	return
}
