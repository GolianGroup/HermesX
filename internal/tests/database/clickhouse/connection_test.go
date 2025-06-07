package clickhouse_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"golang_template/internal/config"
	db "golang_template/internal/database/clickhouse"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func setupTestContainer(t *testing.T) (*config.ClickhouseConfig, func()) {
	ctx := context.TODO()

	req := testcontainers.ContainerRequest{
		Image:        "clickhouse/clickhouse-server:latest",
		ExposedPorts: []string{"9000/tcp"},
		WaitingFor:   wait.ForListeningPort("9000/tcp"),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	require.NoError(t, err)

	mappedPort, err := container.MappedPort(ctx, "9000")
	require.NoError(t, err)

	// First create a config without database specification
	initialCfg := &config.ClickhouseConfig{
		Host:         "localhost",
		Port:         fmt.Sprintf("%d", mappedPort.Int()),
		Database:     "", // Empty database name
		User:         "default",
		Password:     "",
		Debug:        true,
		MaxOpenConns: 10,
		MaxIdleConns: 5,
	}

	// Create initial connection
	database, err := db.NewClickhouseDatabase(ctx, initialCfg)
	require.NoError(t, err)

	// Create the test database
	err = database.Exec(ctx, "CREATE DATABASE IF NOT EXISTS test_db")
	require.NoError(t, err)
	database.Close()

	// Now create the final config with the test database
	cfg := &config.ClickhouseConfig{
		Host:         "localhost",
		Port:         fmt.Sprintf("%d", mappedPort.Int()),
		Database:     "test_db",
		User:         "default",
		Password:     "",
		Debug:        true,
		MaxOpenConns: 10,
		MaxIdleConns: 5,
	}

	cleanup := func() {
		container.Terminate(ctx)
	}

	return cfg, cleanup
}

func createTestTable(t *testing.T, database db.ClickhouseDatabase) {
	ctx := context.Background()

	createTable := `
		CREATE TABLE IF NOT EXISTS test_db.logs (
			event_name String,
			timestamp DateTime,
			value Float64
		) ENGINE = MergeTree()
		ORDER BY (timestamp)
	`
	err := database.Exec(ctx, createTable)
	require.NoError(t, err)
}

func TestClickhouseConnection(t *testing.T) {
	cfg, cleanup := setupTestContainer(t)
	defer cleanup()

	ctx := context.Background()

	t.Run("test connection and ping", func(t *testing.T) {
		database, err := db.NewClickhouseDatabase(ctx, cfg)
		require.NoError(t, err)
		defer database.Close()

		err = database.Ping(ctx)
		assert.NoError(t, err)
	})

	t.Run("test bulk insert and fetch", func(t *testing.T) {
		database, err := db.NewClickhouseDatabase(ctx, cfg)
		require.NoError(t, err)
		defer database.Close()

		createTestTable(t, database)

		// Test data
		now := time.Now().UTC().Round(time.Second)
		testData := [][]interface{}{
			{"event1", now, 1.1},
			{"event2", now, 2.2},
			{"event3", now, 3.3},
		}

		// Test bulk insert
		err = database.BulkInsert(ctx, "test_db.logs", testData)
		assert.NoError(t, err)

		// Test bulk fetch
		query := `
	SELECT event_name, timestamp, value
	FROM test_db.logs
	`
		results, err := database.BulkFetch(ctx, query)
		assert.NoError(t, err)
		assert.Equal(t, len(testData), len(results))

		// Verify each row
		for i, row := range results {
			assert.Equal(t, testData[i][0], row[0].(string))    // event_name is String
			assert.Equal(t, testData[i][1], row[1].(time.Time)) // timestamp is DateTime
			assert.Equal(t, testData[i][2], row[2].(float64))   // value is Float64
		}

		// Verify aggregation query
		aggregateQuery := `
				SELECT
					event_name,
					count() as count,
					sum(value) as total
				FROM test_db.logs
				GROUP BY event_name
				ORDER BY event_name
			`
		rows, err := database.Query(ctx, aggregateQuery)
		assert.NoError(t, err)
		defer rows.Close()

		eventCounts := make(map[string]struct {
			count uint64
			total float64
		})

		for rows.Next() {
			var eventName string
			var count uint64
			var total float64
			err := rows.Scan(&eventName, &count, &total)
			assert.NoError(t, err)
			eventCounts[eventName] = struct {
				count uint64
				total float64
			}{count, total}
		}

		assert.Equal(t, 3, len(eventCounts))
		for eventName, data := range eventCounts {
			assert.Equal(t, data.count, uint64(1), "Event %s should appear once", eventName)
		}
	})

	t.Run("test error cases", func(t *testing.T) {
		database, err := db.NewClickhouseDatabase(ctx, cfg)
		require.NoError(t, err)
		defer database.Close()

		// Test invalid query
		_, err = database.Query(ctx, "SELECT * FROM non_existent_table")
		assert.Error(t, err)

		// Test invalid bulk insert
		err = database.BulkInsert(ctx, "non_existent_table", [][]interface{}{
			{1, "test"},
		})
		assert.Error(t, err)

		// Test invalid bulk fetch
		_, err = database.BulkFetch(ctx, "SELECT * FROM non_existent_table")
		assert.Error(t, err)
	})

	t.Run("test large batch insert", func(t *testing.T) {
		database, err := db.NewClickhouseDatabase(ctx, cfg)
		require.NoError(t, err)
		defer database.Close()

		// Cleanup the table before starting
		err = database.Exec(ctx, "TRUNCATE TABLE test_db.logs")
		require.NoError(t, err)

		createTestTable(t, database)

		// Create large batch of test data
		now := time.Now().Round(time.Second)
		var testData [][]interface{}
		for i := 0; i < 10000; i++ {
			testData = append(testData, []interface{}{
				fmt.Sprintf("event%d", i),
				now.Add(time.Duration(i) * time.Second),
				float64(i) * 1.1,
			})
		}

		// Test bulk insert with large batch
		err = database.BulkInsert(ctx, "test_db.logs", testData)
		fmt.Println(err)
		assert.NoError(t, err)

		// Verify count
		var count uint64
		rows, err := database.Query(ctx, "SELECT count() FROM test_db.logs")
		require.NoError(t, err)
		defer rows.Close()

		require.True(t, rows.Next(), "Expected at least one row")
		err = rows.Scan(&count)
		assert.NoError(t, err)
		assert.Equal(t, uint64(10000), count)
	})
}
