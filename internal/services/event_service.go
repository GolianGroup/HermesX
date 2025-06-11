package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	dto "hermesx/handler/dtos"
	"hermesx/internal/ent"
	"hermesx/internal/producers"
	"hermesx/internal/repositories"

	"time"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type EventService interface {
	CreateEvent(ctx context.Context, event *dto.Event) (*dto.EventResponse, error)
	GetEventByName(ctx context.Context, name string) (*dto.EventResponse, error)
	GetAllEvents(ctx context.Context) ([]*dto.EventResponse, error)
	RemoveEvent(ctx context.Context, name string) error
}

type eventService struct {
	Repo  repositories.EventRepository
	Cache producers.RedisClient
}

func NewEventService(repo repositories.EventRepository, cache producers.RedisClient) EventService {
	return &eventService{Repo: repo, Cache: cache}
}

func (s *eventService) CreateEvent(ctx context.Context, event *dto.Event) (*dto.EventResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	id := uuid.New()

	entEvent := &ent.Event{
		ID:               id,
		Name:             event.Name,
		IsCritical:       event.IsCritical,
		PreferredChannel: event.PreferredChannel,
		Metadata:         event.Metadata,
	}

	exists, err := s.Repo.DoesExist(ctx, event.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to check event existance: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("event with name %s already exists", event.Name)
	}

	err = s.Repo.Create(ctx, entEvent)
	if err != nil {
		return nil, fmt.Errorf("failed to create event: %w", err)
	}

	// Store event in cache
	cacheKey := fmt.Sprintf("event:%s", entEvent.Name)
	eventData, _ := json.Marshal(entEvent)
	s.Cache.Set(ctx, cacheKey, eventData, 10*time.Minute)

	return &dto.EventResponse{
		ID:               entEvent.ID.String(),
		Name:             entEvent.Name,
		IsCritical:       entEvent.IsCritical,
		PreferredChannel: entEvent.PreferredChannel,
		Metadata:         entEvent.Metadata,
		CreatedAt:        timestamppb.Now(),
		UpdatedAt:        timestamppb.Now(),
	}, nil
}

func (s *eventService) GetEventByName(ctx context.Context, name string) (*dto.EventResponse, error) {
	if name == "" {
		return nil, errors.New("event name cannot be empty")
	}

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	cacheKey := fmt.Sprintf("event:%s", name)

	// Check cache first
	cachedData, err := s.Cache.Get(ctx, cacheKey)
	if err == nil && cachedData != "" {
		var cachedEvent dto.EventResponse
		err = json.Unmarshal([]byte(cachedData), &cachedEvent)
		if err == nil {
			return &cachedEvent, nil
		}
	}

	entEvent, err := s.Repo.GetByName(ctx, name)
	if err != nil {
		return nil, err
	}

	response := &dto.EventResponse{
		ID:               entEvent.ID.String(),
		Name:             entEvent.Name,
		IsCritical:       entEvent.IsCritical,
		PreferredChannel: entEvent.PreferredChannel,
		Metadata:         entEvent.Metadata,
		CreatedAt:        timestamppb.New(entEvent.CreatedAt),
		UpdatedAt:        timestamppb.New(entEvent.UpdatedAt),
	}

	// Cache the result
	eventData, err := json.Marshal(response)
	if err == nil {
		s.Cache.Set(ctx, cacheKey, eventData, 10*time.Minute)
	}

	return response, nil
}

func (s *eventService) GetAllEvents(ctx context.Context) ([]*dto.EventResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	events, err := s.Repo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get all events: %w", err)
	}

	var responses []*dto.EventResponse
	for _, entEvent := range events {
		responses = append(responses, &dto.EventResponse{
			ID:               entEvent.ID.String(),
			Name:             entEvent.Name,
			IsCritical:       entEvent.IsCritical,
			PreferredChannel: entEvent.PreferredChannel,
			Metadata:         entEvent.Metadata,
			CreatedAt:        timestamppb.New(entEvent.CreatedAt),
			UpdatedAt:        timestamppb.New(entEvent.UpdatedAt),
		})
	}

	return responses, nil
}

func (s *eventService) RemoveEvent(ctx context.Context, name string) error {
	if name == "" {
		return errors.New("event name cannot be empty")
	}

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// Remove from cache
	cacheKey := fmt.Sprintf("event:%s", name)
	s.Cache.Del(ctx, cacheKey)

	// Remove from repository
	err := s.Repo.Remove(ctx, name)
	if err != nil {
		return fmt.Errorf("failed to remove event: %w", err)
	}

	return nil
}
