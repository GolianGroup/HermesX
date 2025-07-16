package repositories

import (
	"context"
	"hermesx/internal/config"
	"hermesx/internal/database/arango"
	"hermesx/internal/repositories/models"
	"log"

	"github.com/arangodb/go-driver/v2/arangodb"
	"go.uber.org/zap"
)

type ProfileRepository interface {
	GetProfile(ctx context.Context, profileID string, userID string) (*models.Profile, error)
}

type profileRepository struct {
	database             arangodb.Database
	profileCollection    arangodb.Collection
	prefrencesCollection arangodb.Collection
	ctx                  context.Context
	config               *config.Config
	logger               *zap.Logger
}

func NewProfileRepository(db arango.ArangoDB, ctx context.Context, cfg *config.Config, logger *zap.Logger) (ProfileRepository, error) {
	profileCollection, err := db.GetCollection(ctx, "profiles")
	if err != nil {
		log.Println("Failed to get profile collection")
		return nil, err
	}

	prefrencesCollection, err := db.GetCollection(ctx, "preferences")
	if err != nil {
		log.Println("failed to get prefrences collection")
		return nil, err
	}

	database := db.Database(ctx)

	return &profileRepository{
		database:             database,
		profileCollection:    profileCollection,
		prefrencesCollection: prefrencesCollection,
		ctx:                  ctx,
		config:               cfg,
		logger:               logger,
	}, nil
}

func (c *profileRepository) GetProfile(ctx context.Context, profileID string, userID string) (*models.Profile, error) {
	select {
	case <-ctx.Done():
		ErrTimeout.Err = ctx.Err()
		return nil, ErrTimeout
	default:
		var profileObj models.Profile

		query := `
			FOR p IN @@collection
				FILTER p._key == @profileID AND p.user_id == @userID AND p.account_status == True
				LIMIT 1
				RETURN p
		`
		opts := arangodb.QueryOptions{
			BindVars: map[string]interface{}{
				"@collection": c.profileCollection.Name(),
				"profileID":   profileID,
				"userID":      userID,
			},
		}

		cursor, err := c.database.Query(ctx, query, &opts)
		if err != nil {
			c.logger.Error("GetProfile repository failed", zap.Error(err))
			return nil, err
		}
		defer cursor.Close()

		if !cursor.HasMore() {
			return nil, ErrProfileNotFound
		}

		_, err = cursor.ReadDocument(ctx, &profileObj)
		if err != nil {
			c.logger.Error("GetProfile repository failed", zap.Error(err))
			return nil, err
		}
		return &profileObj, nil
	}
}
