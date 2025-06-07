package controllers

import (
	"context"
	"golang_template/internal/logging"
	"golang_template/internal/services"

	"golang_template/handler/dtos"

	rpc_service "golang_template/proto"
)

type RpcServiceController interface {
	SayHello(ctx context.Context, req *rpc_service.HelloRequest) (*rpc_service.HelloReply, error)
}

type rpcServiceController struct {
	rpcServiceService services.RpcServiceService
}
func NewRpcServiceController(service services.RpcServiceService, logger logging.Logger) RpcServiceController {
	return &rpcServiceController{
		rpcServiceService: service,
	}
}

func (c *rpcServiceController) SayHello(ctx context.Context, req *rpc_service.HelloRequest) (*rpc_service.HelloReply, error) {
	requestDTO := dto.ToHelloRequestDTO(req)

	responseDTO, err := c.rpcServiceService.SayHello(ctx, requestDTO)
	if err != nil {
		return nil, err
	}

	return responseDTO.ToHelloReply(), nil
}
