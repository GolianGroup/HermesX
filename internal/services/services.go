package services

import "golang_template/internal/repositories"

type Service interface {
	UserService() UserService
	VideoService() VideoService
	RpcServiceService() RpcServiceService
	SystemService() SystemService
}

type service struct {
	userService  UserService
	videoService VideoService
	rpcServiceService RpcServiceService
	systemService SystemService
}

func NewService(repo repositories.Repository) Service {
	userService := NewUserService(repo.UserRepository())
	videoService := NewVideoService(repo.VideoRepository())
	rpcServiceService := NewRpcServiceService()
	systemService := NewSystemService(repo.SystemRepository())
	return &service{
		userService:  userService,
		videoService: videoService,
		rpcServiceService: rpcServiceService,
		systemService: systemService,
	}
}

func (s *service) UserService() UserService {
	return s.userService
}

func (s *service) VideoService() VideoService {
	return s.videoService
}

func (s *service) RpcServiceService() RpcServiceService {
	return s.rpcServiceService
}

func (s *service) SystemService() SystemService {
	return s.systemService
}