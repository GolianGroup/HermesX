package app

import (
	"golang_template/internal/repositories"
	"golang_template/internal/services"
)

func (a *application) InitServices(repository repositories.Repository) services.Service {
	return services.NewService(repository)
}
