package server_test

import (
	"context"
	"github.com/wuyuanyi135/mvprotos/mvcgi"
	"testing"
)

func TestListNothing(t *testing.T) {
	client := GetClient()
	response, err := client.ListDevices(context.Background(), &mvcgi.ListDevicesRequest{})
	if err != nil {
		t.Fatal(err)
	}
	if response.Devices != nil || response.Controller != nil {
		t.Fatal("Should return nothing")
	}
}

func TestListController(t *testing.T) {
	client := GetClient()
	OpenCtrollerServer(client)
	defer CloseServers(client)

	response, err := client.ListDevices(context.Background(), &mvcgi.ListDevicesRequest{
		ListController: &mvcgi.ListDevicesRequest_Options{
			UseCache: false,
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(response)
}
func TestListCamera(t *testing.T) {
	client := GetClient()
	OpenCameraServer(client)
	defer CloseServers(client)

	response, err := client.ListDevices(context.Background(), &mvcgi.ListDevicesRequest{
		ListCamera: &mvcgi.ListDevicesRequest_Options{
			UseCache: false,
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	t.Log(response)
}
func TestListCameraWithCache(t *testing.T) {
	client := GetClient()
	OpenCameraServer(client)
	defer CloseServers(client)

	response, err := client.ListDevices(context.Background(), &mvcgi.ListDevicesRequest{
		ListCamera: &mvcgi.ListDevicesRequest_Options{
			UseCache: true,
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	t.Log(response)
}