package postgres

import (
	"context"
	"fmt"
	"golang_template/internal/config"
	"golang_template/internal/ent"
	"time"

	dbsql "database/sql"

	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
)

type Database interface {
	Close() error
	EntClient() *ent.Client
	DB() *dbsql.DB
}

type database struct {
	pool     *pgxpool.Pool
	database *dbsql.DB
	client   *ent.Client
}

func NewDatabase(ctx context.Context, dbConfig *config.DatabaseConfig) (Database, error) {

	if dbConfig == nil {
		return nil, fmt.Errorf("database config cannot be nil")
	}

	// Validate config values
	if dbConfig.MaxConns < dbConfig.MinConns {
		return nil, fmt.Errorf("maxConns must be greater than or equal to minConns")
	}

	poolConfig, err := pgxpool.ParseConfig(config.GetDSN(dbConfig))
	if err != nil {
		return nil, err
	}

	// Configure pool
	poolConfig.MaxConns = int32(dbConfig.MaxConns)
	poolConfig.MinConns = int32(dbConfig.MinConns)

	timeoutCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(timeoutCtx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("creating pool: %w", err)
	}

	// Test connection
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("pinging database: %w", err)
	}

	db := stdlib.OpenDB(*pool.Config().ConnConfig)

	// Create ent driver
	driver := sql.OpenDB(dialect.Postgres, db)

	// Create ent client
	client := ent.NewClient(ent.Driver(driver))

	return &database{
		pool:     pool,
		database: db,
		client:   client,
	}, nil
}

func (db *database) Close() error {
	var errs []error

	if err := db.client.Close(); err != nil {
		errs = append(errs, fmt.Errorf("closing ent client: %w", err))
	}

	db.pool.Close()
	if err := db.database.Close(); err != nil {
		errs = append(errs, fmt.Errorf("closing database: %w", err))
	}

	if len(errs) > 0 {
		return fmt.Errorf("closing database: %v", errs)
	}
	return nil
}

func (db *database) EntClient() *ent.Client {
	return db.client
}

func (db *database) DB() *dbsql.DB {
	return db.database
}
