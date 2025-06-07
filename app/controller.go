package app

import (
	"golang_template/handler/controllers"
	"golang_template/internal/services"

	"go.uber.org/zap"
)

func (a *application) InitController(service services.Service, logger *zap.Logger) controllers.Controllers {
	return controllers.NewControllers(service, logger)
}
