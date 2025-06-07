package app

import (
	"golang_template/handler/controllers"
	"golang_template/internal/logging"
	"golang_template/internal/services"
)

func (a *application) InitController(service services.Service, logger logging.Logger) controllers.Controllers {
	return controllers.NewControllers(service, logger)
}

