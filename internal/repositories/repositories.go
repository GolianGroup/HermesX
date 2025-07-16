package repositories

import (
	"context"
	"hermesx/internal/database/arango"
	"hermesx/internal/database/postgres"
	"hermesx/internal/database/scylla"
	"hermesx/internal/producers"

	"go.uber.org/zap"
)

type Repository interface {
	SystemRepository() SystemRepository
	EventRepository() EventRepository
	NotificationRepository() NotificationRepository
}

type repository struct {
	systemRepository       SystemRepository
	eventRepository        EventRepository
	notificationRepository NotificationRepository
}

func NewRepository(db postgres.Database, arango arango.ArangoDB, logger *zap.Logger, redis producers.RedisClient, scylla scylla.ScyllaDB, ctx context.Context) Repository {
	systemRepository := NewSystemRepository(db, arango, redis)
	eventRepository := NewEventRepository(db.EntClient())
	notificationRepository := NewNotificationRepository(scylla, logger)
	return &repository{
		systemRepository:       systemRepository,
		eventRepository:        eventRepository,
		notificationRepository: notificationRepository,
	}
}

func (r *repository) SystemRepository() SystemRepository {
	return r.systemRepository
}

func (r *repository) EventRepository() EventRepository {
	return r.eventRepository
}

func (r *repository) NotificationRepository() NotificationRepository {
	return r.notificationRepository
}
