package routers

import (
	"hermesx/handler/controllers"

	"github.com/gofiber/fiber/v2"
)

type RestEventRouter interface {
	AddRoutes(router fiber.Router)
}

type restEventRouter struct {
	Controller controllers.RestEventController
}

func NewRestEventRouter(eventController controllers.RestEventController) RestEventRouter {
	return &restEventRouter{Controller: eventController}
}

func (r restEventRouter) AddRoutes(router fiber.Router) {
	router.Post("", r.Controller.CreateEvent)
	router.Get("", r.Controller.GetAllEvents)
	router.Get("/:name", r.Controller.GetEventByName)
	router.Delete("/:name", r.Controller.RemoveEvent)
}
