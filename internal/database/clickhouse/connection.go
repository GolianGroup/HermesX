package clickhouse

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"time"

	"golang_template/internal/config"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

type ClickhouseDatabase interface {
	Close() error
	Ping(ctx context.Context) error
	Query(ctx context.Context, query string, args ...interface{}) (driver.Rows, error)
	Exec(ctx context.Context, query string, args ...interface{}) error
	Conn() driver.Conn
	BulkInsert(ctx context.Context, table string, data [][]interface{}) error
	BulkFetch(ctx context.Context, query string, args ...interface{}) ([][]interface{}, error)
}

type clickhouseDatabase struct {
	conn driver.Conn
}

func NewClickhouseDatabase(ctx context.Context, cfg *config.ClickhouseConfig) (ClickhouseDatabase, error) {
	options := &clickhouse.Options{
		Addr: []string{config.GetClickhouseAddr(cfg)},
		Auth: clickhouse.Auth{
			Database: cfg.Database,
			Username: cfg.User,
			Password: cfg.Password,
		},
		Debug:        cfg.Debug,
		MaxOpenConns: cfg.MaxOpenConns,
		MaxIdleConns: cfg.MaxIdleConns,
		Compression: &clickhouse.Compression{
			Method: clickhouse.CompressionLZ4,
		},
	}

	conn, err := clickhouse.Open(options)
	if err != nil {
		return nil, fmt.Errorf("failed to create clickhouse connection: %w", err)
	}

	if err := conn.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping clickhouse: %w", err)
	}

	return &clickhouseDatabase{conn: conn}, nil
}

func (db *clickhouseDatabase) Close() error {
	if db.conn != nil {
		return db.conn.Close()
	}
	return nil
}

func (db *clickhouseDatabase) Ping(ctx context.Context) error {
	return db.conn.Ping(ctx)
}

func (db *clickhouseDatabase) Query(ctx context.Context, query string, args ...interface{}) (driver.Rows, error) {
	return db.conn.Query(ctx, query, args...)
}

func (db *clickhouseDatabase) Exec(ctx context.Context, query string, args ...interface{}) error {
	return db.conn.Exec(ctx, query, args...)
}

func (db *clickhouseDatabase) Conn() driver.Conn {
	return db.conn
}

// BulkInsert inserts multiple rows into the specified table.
func (db *clickhouseDatabase) BulkInsert(ctx context.Context, table string, data [][]interface{}) error {
	if len(data) == 0 {
		return nil
	}

	// Create placeholders for the columns (?, ?, ?, etc.)
	placeholders := make([]string, len(data[0]))
	for i := range placeholders {
		placeholders[i] = "?"
	}

	// Construct the query with proper placeholders
	query := fmt.Sprintf("INSERT INTO %s VALUES (%s)", table, strings.Join(placeholders, ", "))

	// Prepare the batch
	batch, err := db.conn.PrepareBatch(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to prepare batch insert: %w", err)
	}
	defer batch.Abort() // Will be ignored if Send() is successful

	// Add all rows to the batch
	for _, row := range data {
		if err := batch.Append(row...); err != nil {
			return fmt.Errorf("failed to append row to batch: %w", err)
		}
	}

	// Send the batch
	if err := batch.Send(); err != nil {
		return fmt.Errorf("failed to send batch: %w", err)
	}

	return nil
}

// BulkFetch fetches data in bulk using the specified query and arguments.
func (db *clickhouseDatabase) BulkFetch(ctx context.Context, query string, args ...interface{}) ([][]interface{}, error) {
	rows, err := db.conn.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute bulk fetch: %w", err)
	}
	defer rows.Close()

	// Get column types
	columnTypes := rows.ColumnTypes()
	columnCount := len(columnTypes)

	var result [][]interface{}

	// Create a map for known column type handling
	scanFuncs := map[string]func() interface{}{
		"string":    func() interface{} { return new(string) },
		"time.Time": func() interface{} { return new(time.Time) },
		"float64":   func() interface{} { return new(float64) },
	}

	for rows.Next() {
		// Prepare the scan arguments
		scanArgs := make([]interface{}, columnCount)

		// Assign appropriate scan arguments based on column types
		for i, columnType := range columnTypes {
			columnTypeStr := columnType.ScanType().String()

			if scanFunc, ok := scanFuncs[columnTypeStr]; ok {
				scanArgs[i] = scanFunc()
			} else {
				scanArgs[i] = new(interface{}) // Generic fallback
			}
		}

		// Scan the row into the prepared arguments
		if err := rows.Scan(scanArgs...); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		// Extract values from pointers
		row := make([]interface{}, columnCount)
		for i, arg := range scanArgs {
			row[i] = reflect.ValueOf(arg).Elem().Interface()
		}

		result = append(result, row)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during row iteration: %w", err)
	}

	return result, nil
}
