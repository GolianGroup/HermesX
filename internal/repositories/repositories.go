package repositories

import (
	"context"
	"golang_template/internal/database/arango"
	"golang_template/internal/database/postgres"
	"golang_template/internal/producers"
)

type Repository interface {
	SystemRepository() SystemRepository
}

// var (
// ErrGlobal = errors.New("some global error")
// )

type repository struct {
	systemRepository SystemRepository
}

func NewRepository(db postgres.Database, arango arango.ArangoDB, redis producers.RedisClient, ctx context.Context) Repository {
	systemRepository := NewSystemRepository(db, arango, redis)
	return &repository{

		systemRepository: systemRepository,
	}
}

func (r *repository) SystemRepository() SystemRepository {
	return r.systemRepository
}
