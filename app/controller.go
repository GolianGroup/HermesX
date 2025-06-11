package app

import (
	"hermesx/handler/controllers"
	"hermesx/internal/services"

	"go.uber.org/zap"
)

func (a *application) InitController(service services.Service, logger *zap.Logger) controllers.Controllers {
	return controllers.NewControllers(service, logger)
}
