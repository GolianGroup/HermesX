package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"hermesx/internal/config"
	"hermesx/internal/database/scylla"

	"github.com/gocql/gocql"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// scyllaMakemigrationCmd represents the ScyllaDB migration creation command
var scyllaMakemigrationCmd = &cobra.Command{
	Use:   "scylla_makemigration [name]",
	Short: "Create a new CQL migration file with UP and DOWN sections",
	Long: `Example usage:
scylla_makemigration add_users_table   # Creates a new migration file
scylla_makemigration add_users_table --dir ./database/scylla/migrations  # Custom directory`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Println("Running scylla_makemigration...")

		if len(args) < 1 {
			cmd.PrintErrln("Error: Please provide a migration name")
			return
		}

		name := args[0]
		dirFlag, err := cmd.Flags().GetString("dir")
		if err != nil {
			cmd.PrintErrf("Error retrieving dir flag: %v\n", err)
			return
		}

		// Load config
		dbConfig, err := config.LoadConfig("config/config.yml")
		if err != nil {
			log.Fatalf("Failed to load config: %v", err)
			return
		}

		// Initialize logger
		logger, _ := zap.NewProduction()
		defer logger.Sync()

		// Initialize ScyllaDB connection
		scyllaDB, err := scylla.NewScyllaDB(cmd.Context(), dbConfig, logger)
		if err != nil {
			log.Fatalf("Failed to connect to ScyllaDB: %v", err)
			return
		}
		defer scyllaDB.Close()

		session := scyllaDB.Session()
		if session == nil {
			log.Fatalf("Failed to get ScyllaDB session")
			return
		}

		// Create migration file
		err = createMigrationFile(dirFlag, name, dbConfig.ScyllaDB.Keyspace, session)
		if err != nil {
			cmd.PrintErrf("Error creating migration file: %v\n", err)
			return
		}

		cmd.Println("Migration file created successfully")
	},
}

func init() {
	RootCmd.AddCommand(scyllaMakemigrationCmd)
	scyllaMakemigrationCmd.Flags().String("dir", "./internal/database/scylla/migrations", "Directory for ScyllaDB migrations")
}

// createMigrationFile generates a new migration file
func createMigrationFile(directory, name, keyspace string, session *gocql.Session) error {
	// Ensure the migration directory exists
	if err := os.MkdirAll(directory, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Check if the table exists in ScyllaDB
	var tableExists bool
	query := `SELECT table_name FROM system_schema.tables WHERE keyspace_name = ? AND table_name = ?`
	iter := session.Query(query, keyspace, name).Iter()
	tableExists = iter.NumRows() > 0
	iter.Close()

	// Generate timestamp-based filename
	timestamp := time.Now().Format("20060102150405")
	filename := fmt.Sprintf("%s_%s.cql", timestamp, name)
	filePath := filepath.Join(directory, filename)

	// Define migration template
	var content string
	if tableExists {
		content = fmt.Sprintf(`-- Migration: %s
-- Keyspace: %s

-- UP Migration: Alter existing table
ALTER TABLE %s.%s ADD column_name data_type; -- Modify as needed

`, name, keyspace, keyspace, name)
	} else {
		content = fmt.Sprintf(`-- Migration: %s
-- Keyspace: %s

-- UP Migration: Create new table
CREATE TABLE IF NOT EXISTS %s.%s (
    column_name data_type -- Modify as needed
);
`, name, keyspace, keyspace, name)
	}

	// Write to file
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write migration file: %w", err)
	}

	log.Printf("Migration file created: %s", filePath)
	return nil
}
