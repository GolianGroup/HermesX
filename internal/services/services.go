package services

import (
	"hermesx/internal/producers"
	"hermesx/internal/repositories"
)

type Service interface {
	SystemService() SystemService
	EventService() EventService
}

type service struct {
	systemService SystemService
	eventService  EventService
}

func NewService(repo repositories.Repository, redis producers.RedisClient) Service {
	systemService := NewSystemService(repo.SystemRepository())
	eventService := NewEventService(repo.EventRepository(), redis)

	return &service{
		systemService: systemService,
		eventService:  eventService,
	}
}

func (s *service) SystemService() SystemService {
	return s.systemService
}

func (s *service) EventService() EventService {
	return s.eventService
}
