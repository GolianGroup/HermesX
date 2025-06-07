package repositories

import (
	"context"
	"fmt"
	"golang_template/internal/database/arango"
	"golang_template/internal/repositories/models"
	"log"

	"github.com/arangodb/go-driver/v2/arangodb"
)

type VideoRepository interface {
	Get(key string) (*models.Video, error)
	Create(video models.Video) error
	Update(videoUpdate models.Video) (*models.Video, error)
	Delete(key string) error
	GetByName(name string) (*models.Video, error)
}

type videoRepository struct {
	database        arangodb.Database
	videoCollection arangodb.Collection
	ctx             context.Context
}

func NewVideoRepository(db arango.ArangoDB, ctx context.Context) (VideoRepository, error) {
	videoCollection, err := db.GetCollection(ctx, "videos_collection")
	if err != nil {
		log.Println("Failed to get videos collection")
		return nil, err
	}

	database := db.Database(ctx)

	return &videoRepository{
		database:        database,
		videoCollection: videoCollection,
		ctx:             ctx,
	}, nil
}

func (c videoRepository) Get(key string) (*models.Video, error) {
	var video models.Video

	opts := arangodb.CollectionDocumentReadOptions{}
	_, err := c.videoCollection.ReadDocumentWithOptions(c.ctx, key, &video, &opts)
	if err != nil {
		return nil, err
	}

	return &video, nil
}

func (c videoRepository) Create(video models.Video) error {
	opts := arangodb.CollectionDocumentCreateOptions{}
	_, err := c.videoCollection.CreateDocumentWithOptions(c.ctx, video, &opts)
	if err != nil {
		return err
	}
	return nil
}

func (c videoRepository) Update(videoUpdate models.Video) (*models.Video, error) {
	var video models.Video
	withWaitForSync := true
	keepNull := true
	options := arangodb.CollectionDocumentUpdateOptions{
		WithWaitForSync: &withWaitForSync,
		NewObject:       &video,
		KeepNull:        &keepNull,
	}
	_, err := c.videoCollection.UpdateDocumentWithOptions(c.ctx, videoUpdate.Key, videoUpdate, &options)
	if err != nil {
		return nil, err
	}

	return &video, nil
}

func (c videoRepository) Delete(key string) error {
	_, err := c.videoCollection.DeleteDocumentWithOptions(c.ctx, key, nil)
	if err != nil {
		return err
	}
	return nil
}

func (c videoRepository) GetByName(name string) (*models.Video, error) {
	var video models.Video
	query := `
	FOR v IN @@collection
		FILTER v.name == @name
		LIMIT 1
		RETURN v
	`
	opts := arangodb.QueryOptions{
		BindVars: map[string]interface{}{
			"@collection": c.videoCollection.Name(),
			"name":        name,
		},
	}
	cursor, err := c.database.Query(c.ctx, query, &opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close()

	if !cursor.HasMore() {
		return nil, fmt.Errorf("video not found")
	}

	_, err = cursor.ReadDocument(c.ctx, &video)
	if err != nil {
		return nil, err
	}

	return &video, nil
}
