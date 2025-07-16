package controllers

import (
	"hermesx/internal/services"
	rpc_service "hermesx/proto/event"

	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

type Controllers interface {
	SystemController() SystemController
	EventController() rpc_service.EventServiceServer
	NotificationController() NotificationController
}

type controllers struct {
	systemController       SystemController
	eventController        rpc_service.EventServiceServer
	notificationController NotificationController
}

func NewControllers(s services.Service, logger *zap.Logger) Controllers {
	validator := validator.New()
	systemController := NewSystemController(s.SystemService(), logger)
	eventController := NewEventController(s.EventService(), logger, validator)
	notificationController := NewNotificationController(s.NotificationService(), logger)
	return &controllers{
		systemController:       systemController,
		eventController:        eventController,
		notificationController: notificationController,
	}
}

func (c *controllers) SystemController() SystemController {
	return c.systemController
}

func (c *controllers) EventController() rpc_service.EventServiceServer {
	return c.eventController
}

func (c *controllers) NotificationController() NotificationController {
	return c.notificationController
}
