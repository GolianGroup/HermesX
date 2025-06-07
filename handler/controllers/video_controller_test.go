package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	dto "golang_template/handler/dtos"
	"golang_template/internal/mocks"
	"golang_template/internal/repositories/models"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"gotest.tools/v3/assert"
)

func TestVideoController(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockVideoService(ctrl)

	controller := NewVideoController(mockService)

	app := fiber.New()

	t.Run("Get video by key successfully", func(t *testing.T) {
		app.Get("/video/:key", controller.GetVideo)

		var video models.Video
		mockService.EXPECT().GetVideo("123").Return(&video, nil)

		req := httptest.NewRequest("GET", "/video/123", nil)
		resp, err := app.Test(req)

		require.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	})

	t.Run("Get video by key failed", func(t *testing.T) {
		app.Get("/video/:key", controller.GetVideo)

		mockService.EXPECT().GetVideo("123").Return(nil, errors.New("failed to get video"))

		req := httptest.NewRequest("GET", "/video/123", nil)
		resp, err := app.Test(req)

		require.NoError(t, err)
		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	})

	t.Run("Create video successfully", func(t *testing.T) {
		app.Post("/video", controller.CreateVideo)

		request := dto.Video{
			Name:        "name",
			Publishable: true,
			Categories:  []string{"cat1", "cat2"},
			Description: "description",
			Type:        "movie",
		}

		body, err := json.Marshal(request)
		require.NoError(t, err)

		mockService.EXPECT().CreateVideo(request).Return(nil).Times(1)

		req := httptest.NewRequest("POST", "/video", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)

		require.NoError(t, err)
		assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
	})

	t.Run("Create video failed with invalid request", func(t *testing.T) {
		app.Post("/video", controller.CreateVideo)

		req := httptest.NewRequest("POST", "/video", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)

		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Create video failed", func(t *testing.T) {
		app.Post("/video", controller.CreateVideo)

		mockService.EXPECT().CreateVideo(gomock.Any()).Return(errors.New("failed to create video")).Times(1)

		request := dto.Video{
			Name:        "name",
			Publishable: true,
			Categories:  []string{"cat1", "cat2"},
			Description: "description",
			Type:        "movie",
		}

		body, err := json.Marshal(request)
		require.NoError(t, err)

		req := httptest.NewRequest("POST", "/video", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		require.NoError(t, err)

		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	})

	t.Run("Update vide successfully", func(t *testing.T) {
		app.Patch("/video/:key", controller.UpdateVideo)

		request := dto.VideoUpdate{
			Key:         "123",
			Categories:  []string{"cat1", "cat2"},
			Description: "description",
			Name:        "name",
			Views:       1,
		}

		body, err := json.Marshal(request)
		require.NoError(t, err)

		var video models.Video
		mockService.EXPECT().UpdateVideo(request).Return(&video, nil).Times(1)

		req := httptest.NewRequest("PATCH", "/video/123", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		require.NoError(t, err)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	})

	t.Run("Update video failed with invalid request", func(t *testing.T) {
		app.Patch("/video/:key", controller.UpdateVideo)

		req := httptest.NewRequest("PATCH", "/video/123", nil)
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		require.NoError(t, err)

		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Update video failed", func(t *testing.T) {
		app.Patch("/video/:key", controller.UpdateVideo)

		request := dto.VideoUpdate{
			Key:         "123",
			Categories:  []string{"cat1", "cat2"},
			Description: "description",
			Name:        "name",
			Views:       1,
		}

		body, err := json.Marshal(request)
		require.NoError(t, err)

		mockService.EXPECT().UpdateVideo(request).Return(nil, errors.New("failed to update video")).Times(1)

		req := httptest.NewRequest("PATCH", "/video/123", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		require.NoError(t, err)

		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	})

	t.Run("Delete video successfully", func(t *testing.T) {
		app.Delete("/video/:key", controller.DeleteVideo)

		mockService.EXPECT().DeleteVideo("123").Return(nil).Times(1)

		req := httptest.NewRequest("DELETE", "/video/123", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	})

	t.Run("Delete video failed", func(t *testing.T) {
		app.Delete("/video/:key", controller.DeleteVideo)

		mockService.EXPECT().DeleteVideo("123").Return(errors.New("failed to delete video")).Times(1)

		req := httptest.NewRequest("DELETE", "/video/123", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)

		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	})

	t.Run("Get video by name successfully", func(t *testing.T) {
		app.Get("/video", controller.GetVideoByName)

		var video models.Video
		mockService.EXPECT().GetVideoByName("name").Return(&video, nil).Times(1)

		req := httptest.NewRequest("GET", "/video?name=name", nil)
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		require.NoError(t, err)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	})

	t.Run("Failed to get video by name", func(t *testing.T) {
		app.Get("/video", controller.GetVideoByName)

		mockService.EXPECT().GetVideoByName("name").Return(nil, errors.New("failed to get video")).Times(1)

		req := httptest.NewRequest("GET", "/video?name=name", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)

		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	})
}
