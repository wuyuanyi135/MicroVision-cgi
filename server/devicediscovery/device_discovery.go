package devicediscovery

import (
	"context"
	"github.com/wuyuanyi135/mvprotos/mvcam"
	"github.com/wuyuanyi135/mvprotos/mvcamctrl"
	"github.com/wuyuanyi135/mvprotos/mvcgi"
	"google.golang.org/grpc"
	"sync"
)

type DeviceDiscoveryServiceImpl struct {
	cameraServer     mvcam.MicroVisionCameraServiceClient
	controllerServer mvcamctrl.MicroVisionCameraControlServiceClient
}

// make sure the client passed in are already connected to the server.
func NewDeviceDiscoveryServiceImpl(
	cameraServer mvcam.MicroVisionCameraServiceClient,
	controllerServer mvcamctrl.MicroVisionCameraControlServiceClient,
) *DeviceDiscoveryServiceImpl {

	return &DeviceDiscoveryServiceImpl{
		cameraServer:     cameraServer,
		controllerServer: controllerServer,
	}
}

func (s *DeviceDiscoveryServiceImpl) DiscoveryDevices(ctx context.Context, req *mvcgi.DiscoveryDevicesRequest) (resp *mvcgi.DiscoveryDevicesResponse, err error) {
	resp = &mvcgi.DiscoveryDevicesResponse{}
	wg := sync.WaitGroup{}
	if req.DiscoverCamera {
		wg.Add(1)
		go func() {
			getDeviceRequest := &mvcam.GetDevicesRequest{
				UseCache: false,
			}
			response, err := s.cameraServer.GetDevices(ctx, getDeviceRequest, grpc.WaitForReady(true))
			if err != nil {
				return
			}
			for _, v := range response.Devices {
				resp.DiscoveredCamera = append(resp.DiscoveredCamera, &mvcgi.Device{Id: v.Id, Connected: v.Connected, DisplayName: v.Model + v.Id})
			}
			wg.Done()
		}()
	}

	if req.DiscoverController {
		wg.Add(1)
		go func() {
			getSerialDevicesRequest := &mvcamctrl.GetSerialDevicesRequest{
				UseCache: false,
			}
			response, err := s.controllerServer.GetSerialDevices(ctx, getSerialDevicesRequest, grpc.WaitForReady(true))
			if err != nil {
				return
			}
			for _, v := range response.DeviceList {
				resp.DiscoveredController = append(resp.DiscoveredController, &mvcgi.Device{Id: v.Destination, DisplayName: v.Name, Connected: v.Connected})
			}
			wg.Done()
		}()
	}
	wg.Wait()
	return
}
