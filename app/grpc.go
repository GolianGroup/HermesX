package app

import (
	"hermesx/handler/controllers"

	rpc_service "hermesx/proto/event"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func (a *application) InitGRPCServer(controller controllers.Controllers, logger *zap.Logger) *grpc.Server {
	grpcServer := grpc.NewServer()

	rpc_service.RegisterEventServiceServer(grpcServer, controller.EventController())

	if a.config.Environment == "development" {
		reflection.Register(grpcServer)
	}

	return grpcServer
}
