package services

import (
	"hermesx/internal/helper/nats"
	"hermesx/internal/producers"
	"hermesx/internal/repositories"

	"go.uber.org/zap"
)

type Service interface {
	SystemService() SystemService
	EventService() EventService
	NotificationService() NotificationService
}

type service struct {
	systemService       SystemService
	eventService        EventService
	notificationService NotificationService
}

func NewService(repo repositories.Repository, nats nats.NatsConnection, redis producers.RedisClient, logger *zap.Logger) Service {
	systemService := NewSystemService(repo.SystemRepository())
	eventService := NewEventService(repo.EventRepository(), redis)
	notificationService := NewNotificationService(repo.NotificationRepository(), repo.EventRepository(), repo.ProfileRepository(), logger, nats)

	return &service{
		systemService:       systemService,
		eventService:        eventService,
		notificationService: notificationService,
	}
}

func (s *service) SystemService() SystemService {
	return s.systemService
}

func (s *service) EventService() EventService {
	return s.eventService
}

func (s *service) NotificationService() NotificationService {
	return s.notificationService
}
