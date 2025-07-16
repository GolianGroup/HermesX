package controllers

import (
	"hermesx/handler/dtos"
	"hermesx/handler/presenters"
	"hermesx/internal/services"
	"hermesx/pkg"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type NotificationController interface {
	CreateNotification(ctx *fiber.Ctx) error
	BroadcastNotification(ctx *fiber.Ctx) error
	ReadNotification(ctx *fiber.Ctx) error
	ReadBroadcast(ctx *fiber.Ctx) error
}

type notificationController struct {
	service  services.NotificationService
	logger   *zap.Logger
	validate *pkg.Validator
	// fallbackController fiber.Handler
}

func NewNotificationController(service services.NotificationService, logger *zap.Logger) NotificationController {
	logger.With(zap.String("controller", "notificationController"))

	// TODO: fallback controller
	return &notificationController{
		service:  service,
		logger:   logger,
		validate: pkg.NewValidator(),
	}
}

func (c *notificationController) CreateNotification(ctx *fiber.Ctx) error {
	var notificationRequest dtos.Notification

	if err := ctx.BodyParser(&notificationRequest); err != nil {
		c.logger.Debug("Error in CreateNotification body parser", zap.Error(err))
		return err
	}

	err := c.validate.ValidateStruct(notificationRequest)
	if err != nil {
		return err
	}

	if err := c.service.SendNotification(ctx.Context(), &notificationRequest); err != nil {
		return err
	}
	return ctx.Status(fiber.StatusAccepted).JSON(
		presenters.SuccessWithDataResponse(
			fiber.Map{
				"title":   notificationRequest.Title,
				"message": notificationRequest.Message,
			},
		),
	)
}

func (c *notificationController) BroadcastNotification(ctx *fiber.Ctx) error {
	var broadcastNotification dtos.Broadcast

	if err := ctx.BodyParser(&broadcastNotification); err != nil {
		c.logger.Debug("error in BroadcastNotification body parser", zap.Error(err))
		return err
	}

	err := c.validate.ValidateStruct(broadcastNotification)
	if err != nil {
		return err
	}

	if err := c.service.BroadcastNotification(ctx.Context(), &broadcastNotification); err != nil {
		return err
	}
	return ctx.Status(fiber.StatusAccepted).JSON(
		presenters.SuccessWithDataResponse(
			fiber.Map{
				"title":   broadcastNotification.Title,
				"message": broadcastNotification.Message,
			},
		),
	)
}

func (c *notificationController) ReadNotification(ctx *fiber.Ctx) error {
	var readNotification dtos.ReadNotification

	if err := ctx.BodyParser(&readNotification); err != nil {
		c.logger.Debug("Error in ReadNotification body parser", zap.Error(err))
		return err
	}

	err := c.validate.ValidateStruct(readNotification)
	if err != nil {
		return err
	}

	if err := c.service.ReadNotification(ctx.Context(), readNotification.NotificationId, readNotification.ProfileId); err != nil {
		return err
	}

	return ctx.Status(fiber.StatusAccepted).JSON(
		presenters.SuccessWithDataResponse(
			fiber.Map{
				"done": true,
			},
		),
	)
}

func (c *notificationController) ReadBroadcast(ctx *fiber.Ctx) error {
	var read dtos.ReadBroadcast

	if err := ctx.BodyParser(&read); err != nil {
		c.logger.Debug("Error in ReadBrodacast body parser", zap.Error(err))
		return err
	}

	err := c.validate.ValidateStruct(read)
	if err != nil {
		return err
	}

	if err := c.service.ReadBroadcastNotification(ctx.Context(), read.BroadcastId, read.ProfileId); err != nil {
		return err
	}

	return ctx.Status(fiber.StatusAccepted).JSON(
		presenters.SuccessWithDataResponse(
			fiber.Map{
				"done": true,
			},
		),
	)
}
