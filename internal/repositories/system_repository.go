package repositories

import (
	"context"
	"hermesx/internal/database/arango"
	"hermesx/internal/database/postgres"
	"hermesx/internal/producers"
)

// FIXME: add the health check data for nats and other dependencies aswell
type SystemRepository interface {
	DBPing(ctx context.Context) error
	ArangoPing(ctx context.Context) error
	RedisPing(ctx context.Context) error
}

type systemRepository struct {
	db     postgres.Database
	arango arango.ArangoDB
	redis  producers.RedisClient
}

func NewSystemRepository(db postgres.Database, arango arango.ArangoDB, redis producers.RedisClient) SystemRepository {
	return &systemRepository{db: db, arango: arango, redis: redis}
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
	if err := r.redis.Ping(ctx).Err(); err != nil {
		return err
	}
	return nil
}
