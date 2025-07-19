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
	eventRouter        EventRouter
	redisClient        producers.RedisClient
}

func NewRouter(controllers controllers.Controllers, redisClient producers.RedisClient) Router {

	systemRouter := NewSystemRouter(controllers.SystemController())
	notificationRouter := NewNotificationRouter(controllers.NotificationController())
	eventRouter := NewRestEventRouter(controllers.RestEventController())

	return &router{
		systemRouter:       systemRouter,
		notificationRouter: notificationRouter,
		eventRouter:        eventRouter,
		redisClient:        redisClient,
	}
}

func (r router) AddRoutes(router fiber.Router) {
	// -----------------------------------------------------
	// TODO: We need authentication for admin and users here
	// -----------------------------------------------------

	notification := router.Group("/notification")
	event := router.Group("/event")

	r.systemRouter.AddRoutes(router)
	r.notificationRouter.AddRoutes(notification)
	r.eventRouter.AddRoutes(event)
}
