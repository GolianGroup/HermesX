package app

import (
	"hermesx/internal/database/arango"
	"hermesx/internal/database/postgres"
	"hermesx/internal/database/scylla"
	"hermesx/internal/producers"
	"hermesx/internal/repositories"

	"go.uber.org/zap"
)

func (a *application) InitRepositories(db postgres.Database, arango arango.ArangoDB, redis producers.RedisClient, logger *zap.Logger, scylla scylla.ScyllaDB) repositories.Repository {
	return repositories.NewRepository(db, arango, logger, redis, scylla, a.ctx, a.config)
}
