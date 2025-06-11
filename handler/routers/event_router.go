package routers

import (
	"hermesx/handler/controllers"

	"github.com/gofiber/fiber/v2"
)

type EventRouter interface {
	AddRoutes(router fiber.Router)
}

type eventRouter struct {
	Controller controllers.EventController
}

func NewEventRouter(eventController controllers.EventController) EventRouter {
	return &eventRouter{Controller: eventController}
}

func (r eventRouter) AddRoutes(router fiber.Router) {
}
