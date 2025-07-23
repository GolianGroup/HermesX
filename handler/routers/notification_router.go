package routers

import (
	"hermesx/handler/controllers"

	"github.com/gofiber/fiber/v2"
)

type NotificationRouter interface {
	AddRoutes(router fiber.Router)
}

type notificationRouter struct {
	Controller controllers.NotificationController
}

func NewNotificationRouter(notificationController controllers.NotificationController) NotificationRouter {
	return &notificationRouter{Controller: notificationController}
}

func (r notificationRouter) AddRoutes(router fiber.Router) {
	router.Post("", r.Controller.CreateNotification)              //FIXME: Make it private and they must require api key or admin auth
	router.Post("/broadcast", r.Controller.BroadcastNotification) //FIXME: Make it private and they must require api key or admin auth

	router.Post("/read", r.Controller.ReadNotification)
	router.Get("/:profileid", r.Controller.GetUserNotifications)
}
