package controllers

import (
	"golang_template/internal/logging"
	"golang_template/internal/services"
)

type Controllers interface {
	UserController() UserController
	VideoController() VideoController
	RpcServiceController() RpcServiceController
	SystemController() SystemController
}

type controllers struct {
	userController  UserController
	videoController VideoController
	rpcServiceController RpcServiceController
	systemController SystemController
}


func NewControllers(s services.Service, logger logging.Logger) Controllers {
	userController := NewUserController(s.UserService(), logger)
	videoController := NewVideoController(s.VideoService())
	rpcServiceController := NewRpcServiceController(s.RpcServiceService(), logger)
	systemController := NewSystemController(s.SystemService(), logger)
	return &controllers{
		userController:  userController,
		videoController: videoController,
		rpcServiceController: rpcServiceController,
		systemController: systemController,
	}
}

func (c *controllers) UserController() UserController {
	return c.userController
}
	
func (c *controllers) VideoController() VideoController {
	return c.videoController
}

func (c *controllers) RpcServiceController() RpcServiceController {
	return c.rpcServiceController
}

func (c *controllers) SystemController() SystemController {
	return c.systemController
}
