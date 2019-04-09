package server

import "google.golang.org/grpc"

func BuildConnectionCameraServer() *grpc.ClientConn {
	conn, err := grpc.Dial("localhost:5074", grpc.WithInsecure()) // TODO: read from configuration
	if err != nil {
		panic(err)
	}
	return conn
}
func BuildConnectionControllerServer() *grpc.ClientConn {
	conn, err := grpc.Dial("localhost:3050", grpc.WithInsecure()) // TODO: read from configuration
	if err != nil {
		panic(err)
	}
	return conn
}
