package app

import (
	"golang_template/internal/database/arango"
	"golang_template/internal/database/clickhouse"
	"golang_template/internal/database/postgres"
	"golang_template/internal/logging"

	"go.uber.org/zap"
)

func (a *application) InitDatabase(logger logging.Logger) postgres.Database {
	db, err := postgres.NewDatabase(a.ctx, &a.config.DB)
	if err != nil {
		logger.Fatal("Failed to start database", zap.Error(err))

	}
	return db
}

func (a *application) InitClickhouseDatabase(logger logging.Logger) clickhouse.ClickhouseDatabase {
	db, err := clickhouse.NewClickhouseDatabase(a.ctx, &a.config.Clickhouse)
	if err != nil {
		logger.Fatal("failed to setup clickhouse database", zap.Error(err))
	}
	return db
}

func (a *application) InitArangoDB(logger logging.Logger) arango.ArangoDB {
	db, err := arango.NewArangoDB(a.ctx, &a.config.ArangoDB)
	if err != nil {
		logger.Fatal("Failed to start arango database", zap.Error(err))
	}
	return db
}
