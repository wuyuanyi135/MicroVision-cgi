package server_test

import (
	"context"
	"github.com/wuyuanyi135/mvprotos/mvcgi"
	"testing"
)
func TestConnectBoth(t *testing.T) {
	client := GetClient()

	response, err := client.BackendServerInterface(
		context.Background(),
		&mvcgi.BackendServerInterfaceRequest{
			CameraConnection: &mvcgi.BackendServerInterfaceRequest_Connection{
				Url:    "localhost:5074",
				Action: mvcgi.ConnectionAction_CONNECT,
			},
			ControllerConnection: &mvcgi.BackendServerInterfaceRequest_Connection{
				Url:    "localhost:3050",
				Action: mvcgi.ConnectionAction_CONNECT,
			},
		},
	)
	if err != nil {
		t.Error(err)
	}
	if response.CameraServer == mvcgi.ConnectionStatus_CONNECTED &&
		response.ControllerServer == mvcgi.ConnectionStatus_CONNECTED {
	} else {
		t.Error(response.CameraServer, response.ControllerServer)
	}

	response, err = client.BackendServerInterface(
		context.Background(),
		&mvcgi.BackendServerInterfaceRequest{
			ControllerConnection: &mvcgi.BackendServerInterfaceRequest_Connection{
				Action: mvcgi.ConnectionAction_DISCONNECT,
			},
			CameraConnection: &mvcgi.BackendServerInterfaceRequest_Connection{
				Action: mvcgi.ConnectionAction_DISCONNECT,
			},
		},
	)
	if err != nil {
		t.Error(err)
	}
	if response.CameraServer == mvcgi.ConnectionStatus_DISCONNECTED &&
		response.ControllerServer == mvcgi.ConnectionStatus_DISCONNECTED {
	} else {
		t.Error(response.CameraServer, response.ControllerServer)
	}
}
func TestReadStatusOnly(t *testing.T) {
	client := GetClient()

	response1, err := client.BackendServerInterface(context.Background(), &mvcgi.BackendServerInterfaceRequest{})
	if err != nil {
		t.Error(err)
	}
	response2, err := client.BackendServerInterface(context.Background(), &mvcgi.BackendServerInterfaceRequest{
		ControllerConnection: nil,
		CameraConnection:     nil,
	})
	if err != nil {
		t.Error(err)
	}
	response3, err := client.BackendServerInterface(context.Background(), &mvcgi.BackendServerInterfaceRequest{
		ControllerConnection: &mvcgi.BackendServerInterfaceRequest_Connection{
			Action: mvcgi.ConnectionAction_NO_OP,
		},
	})
	if err != nil {
		t.Error(err)
	}

	if response1.ControllerServer == response2.ControllerServer &&
		response1.ControllerServer == response3.ControllerServer &&
		response1.CameraServer == response2.CameraServer &&
		response1.CameraServer == response3.CameraServer {
		t.Log(response1.CameraServer.String())
		t.Log(response2.CameraServer.String())
	}
}
func TestConnectCameraOnly(t *testing.T) {
	client := GetClient()

	response, err := client.BackendServerInterface(
		context.Background(),
		&mvcgi.BackendServerInterfaceRequest{
			CameraConnection: &mvcgi.BackendServerInterfaceRequest_Connection{
				Url:    "localhost:5074",
				Action: mvcgi.ConnectionAction_CONNECT,
			},
		},
	)
	if err != nil {
		t.Error(err)
	}
	if response.CameraServer == mvcgi.ConnectionStatus_CONNECTED &&
		response.ControllerServer == mvcgi.ConnectionStatus_DISCONNECTED {
	} else {
		t.Error(response.CameraServer, response.ControllerServer)
	}

	response, err = client.BackendServerInterface(
		context.Background(),
		&mvcgi.BackendServerInterfaceRequest{
			CameraConnection: &mvcgi.BackendServerInterfaceRequest_Connection{
				Action: mvcgi.ConnectionAction_DISCONNECT,
			},
		},
	)
	if err != nil {
		t.Error(err)
	}
	if response.CameraServer == mvcgi.ConnectionStatus_DISCONNECTED &&
		response.ControllerServer == mvcgi.ConnectionStatus_DISCONNECTED {
	} else {
		t.Error(response.CameraServer, response.ControllerServer)
	}
}
func TestConnectControllerOnly(t *testing.T) {
	client := GetClient()

	response, err := client.BackendServerInterface(
		context.Background(),
		&mvcgi.BackendServerInterfaceRequest{
			ControllerConnection: &mvcgi.BackendServerInterfaceRequest_Connection{
				Url:    "localhost:3050",
				Action: mvcgi.ConnectionAction_CONNECT,
			},
		},
	)
	if err != nil {
		t.Error(err)
	}
	if response.CameraServer == mvcgi.ConnectionStatus_DISCONNECTED &&
		response.ControllerServer == mvcgi.ConnectionStatus_CONNECTED {
	} else {
		t.Error(response.CameraServer, response.ControllerServer)
	}

	response, err = client.BackendServerInterface(
		context.Background(),
		&mvcgi.BackendServerInterfaceRequest{
			ControllerConnection: &mvcgi.BackendServerInterfaceRequest_Connection{
				Action: mvcgi.ConnectionAction_DISCONNECT,
			},
		},
	)
	if err != nil {
		t.Error(err)
	}
	if response.CameraServer == mvcgi.ConnectionStatus_DISCONNECTED &&
		response.ControllerServer == mvcgi.ConnectionStatus_DISCONNECTED {
	} else {
		t.Error(response.CameraServer, response.ControllerServer)
	}
}
func TestDisconnectMultipleTimes(t *testing.T) {
	client := GetClient()

	for i := 0; i < 10; i++ {
		_ , err := client.BackendServerInterface(
			context.Background(),
			&mvcgi.BackendServerInterfaceRequest{
				ControllerConnection: &mvcgi.BackendServerInterfaceRequest_Connection{
					Action: mvcgi.ConnectionAction_DISCONNECT,
				},
			},
		)
		if err != nil {
			t.Fatal(err)
		}
	}

}
func TestConnectWrongAddress(t *testing.T) {
	client := GetClient()

	response, err := client.BackendServerInterface(
		context.Background(),
		&mvcgi.BackendServerInterfaceRequest{
			CameraConnection: &mvcgi.BackendServerInterfaceRequest_Connection{
				Url:    "localhost:5077", // wrong address
				Action: mvcgi.ConnectionAction_CONNECT,
			},
		},
	)
	if err == nil {
		t.Fatal(err)
	} else {
		t.Logf("Got error (expected): %s", err)
	}

	response, err = client.BackendServerInterface(
		context.Background(),
		&mvcgi.BackendServerInterfaceRequest{
			ControllerConnection: &mvcgi.BackendServerInterfaceRequest_Connection{
				Action: mvcgi.ConnectionAction_DISCONNECT,
			},
			CameraConnection: &mvcgi.BackendServerInterfaceRequest_Connection{
				Action: mvcgi.ConnectionAction_DISCONNECT,
			},
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	if response.CameraServer == mvcgi.ConnectionStatus_DISCONNECTED &&
		response.ControllerServer == mvcgi.ConnectionStatus_DISCONNECTED {
	} else {
		t.Error(response.CameraServer, response.ControllerServer)
	}
}
