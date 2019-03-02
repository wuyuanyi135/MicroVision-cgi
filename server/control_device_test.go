package server_test

import (
	"context"
	"github.com/wuyuanyi135/mvprotos/mvcgi"
	"testing"
)

func TestControlDoNothing(t *testing.T) {
	client := GetClient()
	response, err := client.DeviceInterface(context.Background(), &mvcgi.DeviceInterfaceRequest{})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(response)
}