package controllers

import (
	dto "golang_template/handler/dtos"
	"golang_template/handler/presenters"
	"golang_template/internal/services"

	"github.com/gofiber/fiber/v2"
)

type VideoController interface {
	GetVideo(ctx *fiber.Ctx) error
	CreateVideo(ctx *fiber.Ctx) error
	UpdateVideo(ctx *fiber.Ctx) error
	DeleteVideo(ctx *fiber.Ctx) error
	GetVideoByName(ctx *fiber.Ctx) error
}

type videoController struct {
	service services.VideoService
}

func NewVideoController(service services.VideoService) VideoController {
	return &videoController{
		service: service,
	}
}

func (c videoController) GetVideo(ctx *fiber.Ctx) error {
	key := ctx.Params("key")
	video, err := c.service.GetVideo(key)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err})
	}

	presenter := presenters.NewVideoPresenter(video)
	return ctx.Status(fiber.StatusOK).JSON(presenter.Present())
}

func (c videoController) CreateVideo(ctx *fiber.Ctx) error {
	var video dto.Video
	err := ctx.BodyParser(&video)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err})
	}
	err = c.service.CreateVideo(video)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err})
	}
	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "success"})
}

func (c videoController) UpdateVideo(ctx *fiber.Ctx) error {
	var videoUpdate dto.VideoUpdate
	err := ctx.BodyParser(&videoUpdate)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err})
	}
	video, err := c.service.UpdateVideo(videoUpdate)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err})
	}
	return ctx.Status(fiber.StatusOK).JSON(video)
}

func (c videoController) DeleteVideo(ctx *fiber.Ctx) error {
	key := ctx.Params("key")
	err := c.service.DeleteVideo(key)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err})
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "success"})
}

func (c videoController) GetVideoByName(ctx *fiber.Ctx) error {
	name := ctx.Query("name")
	video, err := c.service.GetVideoByName(name)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err})
	}
	return ctx.Status(fiber.StatusOK).JSON(video)
}
