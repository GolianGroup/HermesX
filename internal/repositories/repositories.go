package repositories

import (
	"context"
	"hermesx/internal/database/arango"
	"hermesx/internal/database/postgres"
	"hermesx/internal/producers"
)

type Repository interface {
	SystemRepository() SystemRepository
	EventRepository() EventRepository
}

type repository struct {
	systemRepository SystemRepository
	eventRepository  EventRepository
}

func NewRepository(db postgres.Database, arango arango.ArangoDB, redis producers.RedisClient, ctx context.Context) Repository {
	systemRepository := NewSystemRepository(db, arango, redis)
	eventRepository := NewEventRepository(db.EntClient())
	return &repository{
		systemRepository: systemRepository,
		eventRepository:  eventRepository,
	}
}

func (r *repository) SystemRepository() SystemRepository {
	return r.systemRepository
}

func (r *repository) EventRepository() EventRepository {
	return r.eventRepository
}
