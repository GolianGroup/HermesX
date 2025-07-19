package controllers

import (
	"hermesx/handler/dtos"
	"hermesx/handler/presenters"
	"hermesx/internal/services"
	"hermesx/pkg"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type RestEventController interface {
	CreateEvent(ctx *fiber.Ctx) error
	GetAllEvents(ctx *fiber.Ctx) error
	GetEventByName(ctx *fiber.Ctx) error
	RemoveEvent(ctx *fiber.Ctx) error
}

type restEventController struct {
	service   services.EventService
	logger    *zap.Logger
	validator *pkg.Validator
}

func NewRestEventController(service services.EventService, logger *zap.Logger) RestEventController {
	return &restEventController{service: service, logger: logger, validator: pkg.NewValidator()}
}

func (c *restEventController) CreateEvent(ctx *fiber.Ctx) error {
	var eventRequest dtos.Event

	if err := ctx.BodyParser(&eventRequest); err != nil {
		c.logger.Debug("Error in CreateEvent body parser", zap.Error(err))
		return err
	}

	err := c.validator.ValidateStruct(eventRequest)
	if err != nil {
		return err
	}

	response, err := c.service.CreateEvent(ctx.Context(), &eventRequest)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusAccepted).JSON(
		presenters.SuccessWithDataResponse(response),
	)
}

func (c *restEventController) GetAllEvents(ctx *fiber.Ctx) error {
	event, err := c.service.GetAllEvents(ctx.Context())
	if err != nil {
		c.logger.Error("Error fetching events", zap.Error(err))
		return err
	}
	return ctx.Status(fiber.StatusOK).JSON(
		presenters.EventResponse(event),
	)

}

func (c *restEventController) GetEventByName(ctx *fiber.Ctx) error {
	var event dtos.EventByName
	if err := ctx.ParamsParser(&event); err != nil {
		c.logger.Debug("Error in GetEventByNames params parser", zap.Error(err))
		return err
	}

	err := c.validator.ValidateStruct(event)
	if err != nil {
		return err
	}

	result, err := c.service.GetEventByName(ctx.Context(), event.Name)
	if err != nil {
		return err
	}
	return ctx.Status(fiber.StatusOK).JSON(result)
}

func (c *restEventController) RemoveEvent(ctx *fiber.Ctx) error {
	var event dtos.EventByName
	if err := ctx.ParamsParser(&event); err != nil {
		c.logger.Debug("Error in GetEventByNames params parser", zap.Error(err))
		return err
	}

	err := c.validator.ValidateStruct(event)
	if err != nil {
		return err
	}

	err = c.service.RemoveEvent(ctx.Context(), event.Name)
	if err != nil {
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
