package controllers

import (
	"context"
	"errors"
	"strings"

	dto "hermesx/handler/dtos"

	"hermesx/internal/services"
	rpc_service "hermesx/proto/event"

	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

type EventController interface {
	CreateEvent(ctx context.Context, req *rpc_service.CreateEventRequest) (*rpc_service.EventResponse, error)
	GetEventByName(ctx context.Context, req *rpc_service.GetEventByNameRequest) (*rpc_service.EventResponse, error)
	GetAllEvents(ctx context.Context, req *emptypb.Empty) (*rpc_service.EventResponseList, error)
	RemoveEvent(ctx context.Context, req *rpc_service.GetEventByNameRequest) (*rpc_service.RemoveEventResponse, error)
}

type eventController struct {
	rpc_service.UnimplementedEventServiceServer
	service   services.EventService
	logger    *zap.Logger
	validator *validator.Validate
}

func NewEventController(service services.EventService, logger *zap.Logger, validator *validator.Validate) rpc_service.EventServiceServer {
	return &eventController{service: service, logger: logger, validator: validator}
}

func (c *eventController) CreateEvent(ctx context.Context, req *rpc_service.CreateEventRequest) (*rpc_service.EventResponse, error) {
	eventDTO := dto.Event{
		Name:             req.Name,
		IsCritical:       req.IsCritical,
		IsProtected:      req.IsProtected,
		PreferredChannel: req.PreferredChannel,
		Template:         req.Template,
	}

	err := c.validator.Struct(eventDTO)
	if err != nil {
		var validationErrors []string
		for _, e := range err.(validator.ValidationErrors) {
			validationErrors = append(validationErrors,
				"Field '"+e.Field()+"' failed validation for rule '"+e.Tag()+"'.")
		}
		c.logger.Error("Validation errors occurred", zap.Strings("errors", validationErrors))
		return nil, errors.New("validation failed: " + strings.Join(validationErrors, ", "))
	}

	createdEvent, err := c.service.CreateEvent(
		ctx,
		&eventDTO,
	)
	if err != nil {
		c.logger.Error("Error creating event")
		return nil, err
	}

	return &rpc_service.EventResponse{
		Id:               createdEvent.ID,
		Name:             createdEvent.Name,
		IsCritical:       createdEvent.IsCritical,
		IsProtected:      createdEvent.IsProtected,
		PreferredChannel: createdEvent.PreferredChannel,
		Template:         createdEvent.Template,
		CreatedAt:        createdEvent.CreatedAt,
		UpdatedAt:        createdEvent.UpdatedAt,
	}, nil
}

func (c *eventController) GetEventByName(ctx context.Context, req *rpc_service.GetEventByNameRequest) (*rpc_service.EventResponse, error) {

	event, err := c.service.GetEventByName(ctx, req.Name)
	if err != nil {
		c.logger.Error("Error fetching event by name", zap.Error(err))
		return nil, err
	}

	return &rpc_service.EventResponse{
		Id:               event.ID,
		Name:             event.Name,
		IsCritical:       event.IsCritical,
		IsProtected:      event.IsProtected,
		PreferredChannel: event.PreferredChannel,
		Template:         event.Template,
		CreatedAt:        event.CreatedAt,
		UpdatedAt:        event.UpdatedAt,
	}, nil
}

func (c *eventController) GetAllEvents(ctx context.Context, req *emptypb.Empty) (*rpc_service.EventResponseList, error) {
	events, err := c.service.GetAllEvents(ctx)
	if err != nil {
		c.logger.Error("Error fetching events ", zap.Error(err))
		return nil, err
	}

	var eventResponses []*rpc_service.EventResponse
	for _, event := range events {
		eventResponses = append(eventResponses, &rpc_service.EventResponse{
			Id:               event.ID,
			Name:             event.Name,
			IsCritical:       event.IsCritical,
			IsProtected:      event.IsProtected,
			PreferredChannel: event.PreferredChannel,
			Template:         event.Template,
			CreatedAt:        event.CreatedAt,
			UpdatedAt:        event.UpdatedAt,
		})
	}

	return &rpc_service.EventResponseList{
		Events: eventResponses,
	}, nil
}

func (c *eventController) RemoveEvent(ctx context.Context, req *rpc_service.GetEventByNameRequest) (*rpc_service.RemoveEventResponse, error) {
	err := c.service.RemoveEvent(ctx, req.Name)
	if err != nil {
		c.logger.Error("Error removing event", zap.Error(err))
		return nil, err
	}

	return &rpc_service.RemoveEventResponse{
		Message: "Success",
	}, nil
}
