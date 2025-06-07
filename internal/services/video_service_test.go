package services

import (
	"errors"
	dto "golang_template/handler/dtos"
	"golang_template/internal/mocks"
	"golang_template/internal/repositories/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestVideoService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockVideoRepository(ctrl)

	service := NewVideoService(mockRepo)

	t.Run("Get video by key successfully", func(t *testing.T) {
		foundVideo := &models.Video{
			Key:         "123",
			Name:        "name",
			Views:       1,
			Publishable: true,
			Categories:  []string{"cat1", "cat2"},
			Description: "desc",
			Type:        "type",
		}
		mockRepo.EXPECT().Get("123").Return(foundVideo, nil)

		video, err := service.GetVideo("123")
		require.NoError(t, err)

		assert.IsType(t, &models.Video{}, video)
	})

	t.Run("Failed to get video by key", func(t *testing.T) {
		mockRepo.EXPECT().Get("123").Return(nil, errors.New("failed to get video"))

		_, err := service.GetVideo("123")
		require.Error(t, err)

		assert.Equal(t, err.Error(), "failed to get video")
	})

	t.Run("Create video successfully", func(t *testing.T) {
		mockRepo.EXPECT().Create(gomock.Any()).Return(nil)
		var video dto.Video

		err := service.CreateVideo(video)
		require.NoError(t, err)

		assert.Equal(t, err, nil)
	})

	t.Run("Failed to create video", func(t *testing.T) {
		mockRepo.EXPECT().Create(gomock.Any()).Return(errors.New("failed to create video"))
		var video dto.Video

		err := service.CreateVideo(video)
		require.Error(t, err)

		assert.Equal(t, err.Error(), "failed to create video")
	})

	t.Run("Update video successfully", func(t *testing.T) {
		mockRepo.EXPECT().Update(gomock.Any()).Return(&models.Video{}, nil)
		var videoUpdate dto.VideoUpdate

		video, err := service.UpdateVideo(videoUpdate)
		require.NoError(t, err)

		assert.Equal(t, err, nil)
		assert.IsType(t, &models.Video{}, video)
	})

	t.Run("Failed to update video", func(t *testing.T) {
		mockRepo.EXPECT().Update(gomock.Any()).Return(nil, errors.New("failed to update video"))
		var videoUpdate dto.VideoUpdate

		_, err := service.UpdateVideo(videoUpdate)
		require.Error(t, err)

		assert.Equal(t, err.Error(), "failed to update video")
	})

	t.Run("Delete video successfully", func(t *testing.T) {
		mockRepo.EXPECT().Delete("123").Return(nil)

		err := service.DeleteVideo("123")
		require.NoError(t, err)

		assert.Equal(t, err, nil)
	})

	t.Run("Failed to delete video", func(t *testing.T) {
		mockRepo.EXPECT().Delete("123").Return(errors.New("failed to delete video"))

		err := service.DeleteVideo("123")
		require.Error(t, err)

		assert.Equal(t, err.Error(), "failed to delete video")
	})

	t.Run("Get video by name successfully", func(t *testing.T) {
		mockRepo.EXPECT().GetByName("name").Return(&models.Video{}, nil)

		video, err := service.GetVideoByName("name")
		require.NoError(t, err)

		assert.Equal(t, err, nil)
		assert.IsType(t, &models.Video{}, video)
	})

	t.Run("Failed to get video by name", func(t *testing.T) {
		mockRepo.EXPECT().GetByName("name").Return(nil, errors.New("failed to get video"))

		_, err := service.GetVideoByName("name")
		require.Error(t, err)

		assert.Equal(t, err.Error(), "failed to get video")
	})
}
