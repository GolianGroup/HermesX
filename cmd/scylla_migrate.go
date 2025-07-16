package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"

	"hermesx/internal/config"
	"hermesx/internal/database/scylla"

	"github.com/gocql/gocql"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// scyllaMigrateCmd represents the ScyllaDB migration apply command
var scyllaMigrateCmd = &cobra.Command{
	Use:   "scylla_migrate",
	Short: "Apply all pending CQL migration files to ScyllaDB",
	Long: `Example usage:
scylla_migrate   # Applies all pending migrations
scylla_migrate --dir ./database/scylla/migrations  # Custom directory`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Println("Running scylla_migrate...")

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

		// Apply migrations
		err = applyMigrations(dirFlag, session)
		if err != nil {
			cmd.PrintErrf("Error applying migrations: %v\n", err)
			return
		}

		cmd.Println("All migrations applied successfully")
	},
}

func init() {
	RootCmd.AddCommand(scyllaMigrateCmd)
	scyllaMigrateCmd.Flags().String("dir", "./internal/database/scylla/migrations", "Directory for ScyllaDB migrations")
}

// applyMigrations reads migration files and executes them
func applyMigrations(directory string, session *gocql.Session) error {
	// Read migration files from the directory
	files, err := os.ReadDir(directory)
	if err != nil {
		return fmt.Errorf("failed to read migration directory: %w", err)
	}

	// Sort files by timestamp to ensure they are applied in the correct order
	var migrationFiles []string
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".cql" {
			migrationFiles = append(migrationFiles, file.Name())
		}
	}

	// Sort files based on the timestamp in the filename
	sort.Strings(migrationFiles)

	// Loop through each migration file and apply it
	for _, migrationFile := range migrationFiles {
		filePath := filepath.Join(directory, migrationFile)

		// Read the content of the migration file
		content, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", filePath, err)
		}

		// Apply the migration (execute the CQL)
		err = session.Query(string(content)).Exec()
		if err != nil {
			return fmt.Errorf("failed to apply migration %s: %w", migrationFile, err)
		}

		log.Printf("Applied migration: %s", migrationFile)
	}

	return nil
}
