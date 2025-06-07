package repositories

import (
	"context"
	"golang_template/internal/database/arango"
	"golang_template/internal/database/clickhouse"
	"golang_template/internal/database/postgres"
	"golang_template/internal/producers"
)

type SystemRepository interface {
	DBPing(ctx context.Context) error
	ArangoPing(ctx context.Context) error
	RedisPing(ctx context.Context) error
	ClickhousePing(ctx context.Context) error
}

type systemRepository struct {
	db postgres.Database
	arango arango.ArangoDB
	redis producers.RedisClient
	clickhouse clickhouse.ClickhouseDatabase
}

func NewSystemRepository(db postgres.Database, arango arango.ArangoDB, redis producers.RedisClient, clickhouse clickhouse.ClickhouseDatabase) SystemRepository {
	return &systemRepository{db: db, arango: arango, redis: redis, clickhouse: clickhouse}
}


func (r *systemRepository) DBPing(ctx context.Context) error {	
	if err := r.db.DB().Ping(); err != nil {
		return err
	}
	return nil
}

func (r *systemRepository) ArangoPing(ctx context.Context) error {
	if err := r.arango.Ping(ctx); err != nil {
		return err
	}
	return nil
}

func (r *systemRepository) RedisPing(ctx context.Context) error {
	if err := r.redis.RedisStorage().Conn().Ping(ctx).Err(); err != nil {
		return err
	}
	return nil
}

func (r *systemRepository) ClickhousePing(ctx context.Context) error {
	if err := r.clickhouse.Ping(ctx); err != nil {
		return err
	}
	return nil
}

