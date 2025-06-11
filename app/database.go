package app

import (
	"hermesx/internal/database/arango"
	"hermesx/internal/database/postgres"

	"go.uber.org/zap"
)

func (a *application) InitDatabase(logger *zap.Logger) postgres.Database {
	db, err := postgres.NewDatabase(a.ctx, &a.config.DB)
	if err != nil {
		logger.Fatal("Failed to start database", zap.Error(err))

	}
	return db
}

func (a *application) InitArangoDB(logger *zap.Logger) arango.ArangoDB {
	db, err := arango.NewArangoDB(a.ctx, &a.config.ArangoDB)
	if err != nil {
		logger.Fatal("Failed to start arango database", zap.Error(err))
	}
	return db
}
