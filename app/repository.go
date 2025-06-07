package app

import (
	"golang_template/internal/database/arango"
	"golang_template/internal/database/clickhouse"
	"golang_template/internal/database/postgres"
	"golang_template/internal/producers"
	"golang_template/internal/repositories"
)

func (a *application) InitRepositories(db postgres.Database, arango arango.ArangoDB, redis producers.RedisClient, clickhouse clickhouse.ClickhouseDatabase) repositories.Repository {
	return repositories.NewRepository(db, arango, redis, clickhouse, a.ctx)
}
