package routers

import (
	"golang_template/handler/controllers"
	"golang_template/internal/producers"

	"github.com/gofiber/fiber/v2"
)

type UserRouter interface {
	AddRoutes(router fiber.Router)
}

type userRouter struct {
	Controller  controllers.UserController
	redisClient producers.RedisClient
}

func NewUserRouter(userController controllers.UserController, redisClient producers.RedisClient) UserRouter {
	return &userRouter{Controller: userController, redisClient: redisClient}
}

func (r userRouter) AddRoutes(router fiber.Router) {
	// init routes for user
	// has controller
	router.Post("/user", r.Controller.Login)
}
