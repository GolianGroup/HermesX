package repositories

import (
	"context"
	"golang_template/internal/database/arango"
	"golang_template/internal/database/clickhouse"
	"golang_template/internal/database/postgres"
	"golang_template/internal/producers"
	"log"
)

type Repository interface {
	UserRepository() UserRepository
	VideoRepository() VideoRepository
	SystemRepository() SystemRepository
}

// var (
// ErrGlobal = errors.New("some global error")
// )

type repository struct {
	userRepository  UserRepository
	videoRepository VideoRepository
	systemRepository SystemRepository
}

func NewRepository(db postgres.Database, arango arango.ArangoDB, redis producers.RedisClient, clickhouse clickhouse.ClickhouseDatabase, ctx context.Context) Repository {
	userRepository := NewUserRepository(db)
	videoRepository, err := NewVideoRepository(arango, ctx)
	if err != nil {
		log.Panic("Failed to initialize video repository", err)
	}
	systemRepository := NewSystemRepository(db, arango, redis, clickhouse)
	return &repository{
		userRepository:  userRepository,
		videoRepository: videoRepository,
		systemRepository: systemRepository,
	}
}

func (r *repository) UserRepository() UserRepository {
	return r.userRepository
}

func (r *repository) VideoRepository() VideoRepository {
	return r.videoRepository
}

func (r *repository) SystemRepository() SystemRepository {
	return r.systemRepository
}
