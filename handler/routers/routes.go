package routers

import (
	"hermesx/handler/controllers"
	"hermesx/internal/producers"

	"github.com/gofiber/fiber/v2"
)

type Router interface {
	AddRoutes(router fiber.Router)
}

type router struct {
	systemRouter       SystemRouter
	notificationRouter NotificationRouter
	redisClient        producers.RedisClient
}

func NewRouter(controllers controllers.Controllers, redisClient producers.RedisClient) Router {

	systemRouter := NewSystemRouter(controllers.SystemController())
	notificationRouter := NewNotificationRouter(controllers.NotificationController())

	return &router{
		systemRouter:       systemRouter,
		notificationRouter: notificationRouter,
		redisClient:        redisClient,
	}
}

func (r router) AddRoutes(router fiber.Router) {

	notification := router.Group("/notification")

	r.systemRouter.AddRoutes(router)
	r.notificationRouter.AddRoutes(notification)
}
