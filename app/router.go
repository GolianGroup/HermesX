package app

import (
	"hermesx/handler/controllers"
	"hermesx/handler/routers"
	"hermesx/internal/producers"

	"github.com/gofiber/fiber/v2"
)

func (a *application) InitRouter(app *fiber.App, controller controllers.Controllers, redisClient producers.RedisClient) routers.Router {
	router := routers.NewRouter(controller, redisClient)
	router.AddRoutes(app.Group(""))
	return router
}
