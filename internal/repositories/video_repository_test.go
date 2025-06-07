package repositories

import (
	"context"
	"golang_template/internal/mocks"
	"golang_template/internal/repositories/models"
	"testing"

	"github.com/arangodb/go-driver/v2/arangodb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestVideoRepository(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockArango := mocks.NewMockArangoDB(ctrl)
	mockCollection := mocks.NewMockCollection(ctrl)
	mockCollection.EXPECT().Name().Return("videos_collection").AnyTimes()
	mockdb := mocks.NewMockDatabase(ctrl)

	ctx := context.Background()

	mockArango.EXPECT().GetCollection(ctx, "videos_collection").Return(mockCollection, nil).AnyTimes()
	mockArango.EXPECT().Database(ctx).Return(mockdb)

	repo, err := NewVideoRepository(mockArango, ctx)
	require.NoError(t, err)

	t.Run("Get video by key successfully", func(t *testing.T) {
		var meta arangodb.DocumentMeta
		mockCollection.EXPECT().ReadDocumentWithOptions(ctx, "123", gomock.Any(), gomock.Any()).Return(meta, nil)

		video, err := repo.Get("123")
		require.NoError(t, err)

		assert.IsType(t, &models.Video{}, video)
	})

	t.Run("Create video successfully", func(t *testing.T) {
		var response arangodb.CollectionDocumentCreateResponse
		mockCollection.EXPECT().CreateDocumentWithOptions(ctx, gomock.Any(), gomock.Any()).Return(response, nil)
		video := models.Video{
			Name:        "name",
			Publishable: true,
			Categories:  []string{"cat1", "cat2"},
			Description: "desc",
			Type:        "type",
		}
		err := repo.Create(video)
		require.NoError(t, err)
		assert.Equal(t, err, nil)
	})

	t.Run("Update video successfully", func(t *testing.T) {
		var response arangodb.CollectionDocumentUpdateResponse
		mockCollection.EXPECT().UpdateDocumentWithOptions(ctx, "123", gomock.Any(), gomock.Any()).Return(response, nil)
		videoUpdate := models.Video{
			Key:         "123",
			Categories:  []string{"cat1", "cat2"},
			Description: "desc",
			Name:        "name",
			Views:       1,
		}
		video, err := repo.Update(videoUpdate)
		require.NoError(t, err)

		assert.IsType(t, &models.Video{}, video)
	})

	t.Run("Delete video successfully", func(t *testing.T) {
		var response arangodb.CollectionDocumentDeleteResponse
		mockCollection.EXPECT().DeleteDocumentWithOptions(ctx, "123", gomock.Any()).Return(response, nil)
		err := repo.Delete("123")
		require.NoError(t, err)

		assert.Equal(t, err, nil)
	})
}
