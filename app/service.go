package app

import (
	"hermesx/internal/producers"
	"hermesx/internal/repositories"
	"hermesx/internal/services"
)

func (a *application) InitServices(repository repositories.Repository, redis producers.RedisClient) services.Service {
	return services.NewService(repository, redis)
}
