package routers

import (
	"golang_template/handler/controllers"

	"github.com/gofiber/fiber/v2"
)

type VideoRouter interface {
	AddRoutes(router fiber.Router)
}

type videoRouter struct {
	Controller controllers.VideoController
}

func NewVideoRouter(videoController controllers.VideoController) VideoRouter {
	return &videoRouter{Controller: videoController}
}

func (r videoRouter) AddRoutes(router fiber.Router) {
	router.Get("/video/:key", r.Controller.GetVideo)
	router.Post("/video", r.Controller.CreateVideo)
	router.Patch("/video", r.Controller.UpdateVideo)
	router.Delete("/video/:key", r.Controller.DeleteVideo)
	router.Get("/video", r.Controller.GetVideoByName)
}
