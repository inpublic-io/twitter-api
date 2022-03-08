package grpc

import (
	twitterv1 "github.com/inpublic-io/inpublicapis/twitter/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func NewServer() *grpc.Server {
	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)

	twitterv1.RegisterUsersServer(grpcServer, &usersServer{})
	twitterv1.RegisterMetricsServer(grpcServer, &metricsServer{})
	return grpcServer
}
