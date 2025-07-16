package app

import (
	"hermesx/internal/helper/nats"
	"hermesx/internal/producers"
	"hermesx/internal/repositories"
	"hermesx/internal/services"

	"go.uber.org/zap"
)

func (a *application) InitServices(repository repositories.Repository, nats nats.NatsConnection, redis producers.RedisClient, logger *zap.Logger) services.Service {
	return services.NewService(repository, nats, redis, logger)
}
