/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"hermesx/internal/config"
	"log"
	"time"

	"github.com/arangodb/go-driver/v2/arangodb"
	"github.com/arangodb/go-driver/v2/connection"
	"github.com/spf13/cobra"
)

// arangoCreateCmd represents the arangoCreate command
var arangoCreateCmd = &cobra.Command{
	Use:   "ag_create",
	Short: "Creates arangodb database.",
	Long: `ag_create creates database based on environment variables
	ag_create --db_name=mydb creates database with name mydb`, //TODO: better examples
	Run: func(cmd *cobra.Command, args []string) { //TODO: use RunE and return error do not panic
		dbConfig, err := config.LoadConfig("config/config.yml") //TODO: use config flag and use this as fallback the flag must be bound
		if err != nil {
			log.Panicf("failed to setup viper: %s", err) //TODO: use better wording, viper is a package the failure lies within config not viper package
			return
		}

		connStrs, err := config.GetArangoStrings(&dbConfig.ArangoDB)
		if err != nil {
			cmd.PrintErrf("failed to setup arango to create database: %s", err)
		}
		endpoint := connection.NewRoundRobinEndpoints(connStrs)                                                                                                   // TODO: use a function for creating connections
		conn := connection.NewHttp2Connection(connection.DefaultHTTP2ConfigurationWrapper(endpoint /*InsecureSkipVerify*/, dbConfig.ArangoDB.InsecureSkipVerify)) // TODO: remove unneccessary coment

		auth := connection.NewBasicAuth(dbConfig.ArangoDB.User, dbConfig.ArangoDB.Pass)
		err = conn.SetAuthentication(auth)
		if err != nil {
			cmd.PrintErrf("failed to authenticate arango to create database: %s", err)
			return
		}

		client := arangodb.NewClient(conn)

		dbFlag, err := cmd.Flags().GetString("db_name")
		if err != nil {
			cmd.PrintErrf("failed to get database name: %s", err) //TODO: wrap erro
			return
		}

		var dbName string
		if dbFlag == "" {
			dbName = dbConfig.ArangoDB.DBName
		} else {
			dbName = dbFlag
		}
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		dbExists, err := client.DatabaseExists(ctx, dbName)
		if err != nil {
			cmd.PrintErrf("failed to check if database exists: %s", err) //TODO: wrap erro
			return
		}
		if !dbExists {
			_, err = client.CreateDatabase(ctx, dbName, nil)
			if err != nil {
				cmd.PrintErrf("failed to create database: %s", err) //TODO: wrap erro
				return
			}
			cmd.Printf("Database %s created successfully.\n", dbName)
			return
		} else {
			cmd.Printf("Database %s already exists.\n", dbName)
			return
		}
	},
}

func init() {
	RootCmd.AddCommand(arangoCreateCmd)
	arangoCreateCmd.Flags().String("db_name", "", "Database name") // TODO: add verbosity flag
}
