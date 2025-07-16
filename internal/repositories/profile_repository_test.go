package repositories

import (
	"context"
	"hermesx/internal/config"
	"hermesx/internal/mocks"
	"hermesx/internal/repositories/models"
	"testing"
	"time"

	"github.com/arangodb/go-driver/v2/arangodb"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

func TestProfileRepository(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockArango := mocks.NewMockArangoDB(ctrl)
	mockCollectionprofile := mocks.NewMockCollection(ctrl)
	mockCollectionprofile.EXPECT().Name().Return("profiles").AnyTimes()

	mockCollectionprefrences := mocks.NewMockCollection(ctrl)
	mockCollectionprefrences.EXPECT().Name().Return("preferences").AnyTimes()
	mockdb := mocks.NewMockDatabase(ctrl)
	mockCursor := mocks.NewMockCursor(ctrl)

	ctx := context.Background()

	testConfig := &config.Config{}
	logger := zap.NewExample()

	mockArango.EXPECT().GetCollection(ctx, "profiles").Return(mockCollectionprofile, nil).AnyTimes()
	mockArango.EXPECT().Database(ctx).Return(mockdb)

	mockArango.EXPECT().GetCollection(ctx, "preferences").Return(mockCollectionprefrences, nil).AnyTimes()

	repo, err := NewProfileRepository(mockArango, ctx, testConfig, logger)
	require.NoError(t, err)

	t.Run("Get profile not found", func(t *testing.T) {
		profileID := uuid.New().String()
		userID := uuid.New().String()

		mockdb.EXPECT().Query(ctx, gomock.Any(), gomock.Any()).Return(mockCursor, nil)
		mockCursor.EXPECT().HasMore().Return(false)
		mockCursor.EXPECT().Close().Return(nil)

		profile, err := repo.GetProfile(ctx, profileID, userID)
		require.Error(t, err)
		assert.Equal(t, err, ErrProfileNotFound)
		assert.Nil(t, profile)
	})

	t.Run("Get profile successfully", func(t *testing.T) {
		profileID := uuid.New().String()
		userID := uuid.New().String()

		expectedProfile := models.Profile{
			ProfileID:   profileID,
			UserID:      userID,
			ProfileName: "testuser",
			CreatedAt:   time.Now().UTC().Format(time.RFC3339),
		}

		mockdb.EXPECT().Query(ctx, gomock.Any(), gomock.Any()).Return(mockCursor, nil)
		mockCursor.EXPECT().HasMore().Return(true)
		mockCursor.EXPECT().ReadDocument(ctx, gomock.Any()).SetArg(1, expectedProfile).Return(arangodb.DocumentMeta{}, nil)
		mockCursor.EXPECT().Close().Return(nil)

		profile, err := repo.GetProfile(ctx, userID, profileID)
		require.NoError(t, err)
		assert.Equal(t, &expectedProfile, profile)
	})

}
