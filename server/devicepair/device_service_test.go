package devicepair_test

import (
	"context"
	"fmt"
	"net"
	"testing"

	"github.com/golang/protobuf/ptypes"
	"github.com/wuyuanyi135/MicroVisionCGI/server/devicepair"
	"github.com/wuyuanyi135/mvprotos/mvcgi"
	"google.golang.org/grpc"
)

const Port = 30500

var client mvcgi.DevicePairServiceClient

func TestMain(m *testing.M) {
	go StartTestServer()
	conn, err := grpc.Dial(fmt.Sprintf("0.0.0.0:%d", Port), grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	client = mvcgi.NewDevicePairServiceClient(conn)

	m.Run()
}

func StartTestServer() {
	grpcServer := grpc.NewServer()
	devicePairService := devicepair.NewDeviceServiceImpl()
	mvcgi.RegisterDevicePairServiceServer(grpcServer, devicePairService)
	address := fmt.Sprintf("0.0.0.0:%d", Port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		panic(err)
	}
	err = grpcServer.Serve(listener)
	if err != nil {
		panic(err)
	}

}

func TestDeviceServiceImpl_Create(t *testing.T) {
	_, err := client.Create(context.Background(), &mvcgi.CreateDevicePairRequest{
		Device: &mvcgi.DevicePair{
			CameraId:     "test_camera_id",
			ControllerId: "test_controller_id",
			CreatedAt:    ptypes.TimestampNow(),
		},
	}, grpc.WaitForReady(true))

	if err != nil {
		t.Error(err)
		return
	}
}

func TestDeviceServiceImpl_List(t *testing.T) {
	response, err := client.List(context.Background(), &mvcgi.ListDevicePairRequest{}, grpc.WaitForReady(true))

	if err != nil {
		t.Error(err)
		return
	}

	for key, value := range response.Devices {
		t.Logf("%d: camera: %s; controller: %s; created at: %s", key, value.CameraId, value.ControllerId, value.CreatedAt.String())
	}
}
func TestDeviceServiceImpl_Update(t *testing.T) {
	response, err := client.List(context.Background(), &mvcgi.ListDevicePairRequest{}, grpc.WaitForReady(true))

	if err != nil {
		t.Error(err)
		return
	}

	if len(response.Devices) == 0 {
		t.Skip("No item to update, skip.")
		return
	}

	// update with id
	pair := response.Devices[0]
	pair.ControllerId = "updated_controller_id"
	pair.CameraId = "updated_camera_id"
	_, err = client.Update(
		context.Background(),
		&mvcgi.UpdateDevicePairRequest{
			Device:   &mvcgi.UpdateDevicePairRequest_Id{Id: pair.Id},
			NewValue: pair,
		}, grpc.WaitForReady(true))
	if err != nil {
		t.Error(err)
	}
	response, err = client.List(context.Background(), &mvcgi.ListDevicePairRequest{}, grpc.WaitForReady(true))
	if err != nil {
		t.Error(err)
	}
	if response.Devices[0].ControllerId != pair.ControllerId || response.Devices[0].CameraId != pair.CameraId {
		t.Error("Update failed.")
	}

	// update with item
	pair.ControllerId = "updated_controller_id_item"
	pair.CameraId = "updated_camera_id_item"
	_, err = client.Update(
		context.Background(),
		&mvcgi.UpdateDevicePairRequest{
			Device:   &mvcgi.UpdateDevicePairRequest_DevicePair{DevicePair: pair},
			NewValue: pair,
		}, grpc.WaitForReady(true))
	if err != nil {
		t.Error(err)
	}
	response, err = client.List(context.Background(), &mvcgi.ListDevicePairRequest{}, grpc.WaitForReady(true))
	if err != nil {
		t.Error(err)
	}
	if response.Devices[0].ControllerId != pair.ControllerId || response.Devices[0].CameraId != pair.CameraId {
		t.Error("Update failed.")
	}
}

func TestDeviceServiceImpl_Delete(t *testing.T) {
	response, err := client.List(context.Background(), &mvcgi.ListDevicePairRequest{}, grpc.WaitForReady(true))

	if err != nil {
		t.Error(err)
		return
	}

	if len(response.Devices) == 0 {
		t.Skip("No item to delete, skip.")
		return
	}
	_, err = client.Delete(
		context.Background(),
		&mvcgi.DeleteDevicePairRequest{
			Device: &mvcgi.DeleteDevicePairRequest_DevicePair{DevicePair: response.Devices[0]},
		},
		grpc.WaitForReady(true))
	if err != nil {
		t.Error(err)
	}
}
