package devcon

import (
	"context"
	"github.com/hashicorp/go-multierror"
	"github.com/wuyuanyi135/mvprotos/mvcam"
	"github.com/wuyuanyi135/mvprotos/mvcamctrl"
	"github.com/wuyuanyi135/mvprotos/mvcgi"
	"google.golang.org/grpc"
	"sync"
)

func (s *DeviceConnectionServiceImpl) DiscoveryDevices(ctx context.Context, req *mvcgi.DiscoveryDevicesRequest) (resp *mvcgi.DiscoveryDevicesResponse, err error) {
	resp = &mvcgi.DiscoveryDevicesResponse{
		DiscoveredCamera:     nil,
		DiscoveredController: nil,
	}
	if req.UseCache {
		resp = &s.cache
		return
	}
	e := &multierror.Error{}
	wg := sync.WaitGroup{}
	if req.DiscoverCamera {
		wg.Add(1)
		go func() {
			defer wg.Done()
			getDeviceRequest := &mvcam.GetDevicesRequest{
				UseCache: false,
			}
			response, err := s.cameraServer.GetDevices(ctx, getDeviceRequest, grpc.WaitForReady(true))
			if err != nil {
				e = multierror.Append(e, err)
				return
			}
			for _, v := range response.Devices {
				resp.DiscoveredCamera = append(resp.DiscoveredCamera, &mvcgi.Device{Id: v.Id, Connected: v.Connected, DisplayName: v.Model + v.Id})
			}
		}()
	}

	if req.DiscoverController {
		wg.Add(1)
		go func() {
			defer wg.Done()
			getSerialDevicesRequest := &mvcamctrl.GetSerialDevicesRequest{
				UseCache: false,
			}
			response, err := s.controllerServer.GetSerialDevices(ctx, getSerialDevicesRequest, grpc.WaitForReady(true))
			if err != nil {
				e = multierror.Append(e, err)
				return
			}
			for _, v := range response.DeviceList {
				resp.DiscoveredController = append(resp.DiscoveredController, &mvcgi.Device{Id: v.Destination, DisplayName: v.Name, Connected: v.Connected})
			}
		}()
	}
	wg.Wait()
	if e.ErrorOrNil() != nil {
		err = e
		return
	}
	s.cache = *resp
	return
}
