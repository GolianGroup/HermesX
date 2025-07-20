package repositories

import (
	"context"
	"hermesx/internal/config"
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
	ProfileRepository() ProfileRepository
}

type repository struct {
	systemRepository       SystemRepository
	eventRepository        EventRepository
	notificationRepository NotificationRepository
	profileRepository      ProfileRepository
}

func NewRepository(db postgres.Database, arango arango.ArangoDB, logger *zap.Logger, redis producers.RedisClient, scylla scylla.ScyllaDB, ctx context.Context, cfg *config.Config) Repository {
	systemRepository := NewSystemRepository(db, arango, redis)
	eventRepository := NewEventRepository(db.EntClient())
	notificationRepository := NewNotificationRepository(scylla, logger)
	profileRepository, err := NewProfileRepository(arango, ctx, cfg, logger)
	if err != nil {
		logger.Error("Failed to initialize profile repository", zap.Error(err))
	}
	return &repository{
		systemRepository:       systemRepository,
		eventRepository:        eventRepository,
		notificationRepository: notificationRepository,
		profileRepository:      profileRepository,
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

func (r *repository) ProfileRepository() ProfileRepository {
	return r.profileRepository
}
