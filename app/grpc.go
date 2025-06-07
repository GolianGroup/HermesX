package app

import (
	"golang_template/handler/controllers"
	"golang_template/internal/logging"
	rpc_service "golang_template/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func (a *application) InitGRPCServer(controller controllers.Controllers, logger logging.Logger) *grpc.Server {
	grpcServer := grpc.NewServer()
	// Register server with controller
	rpc_service.RegisterRpcServiceServer(grpcServer, controller.RpcServiceController())

	if a.config.Environment == "development" {
		reflection.Register(grpcServer)
	}

	return grpcServer
}
