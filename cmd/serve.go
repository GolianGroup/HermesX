package cmd

import (
	"context"
	"golang_template/app"
	"golang_template/internal/config"
	"log"

	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serve API",
	Run:   serve,
	Args:  cobra.MaximumNArgs(2),
}

func init() {
	RootCmd.AddCommand(serveCmd)
}

func serve(cmd *cobra.Command, args []string) {
	log.Println("serve")
	//viper
	config, err := config.LoadConfig("config/config.yml")
	if err != nil {
		log.Fatalf("failed to setup viper: %s", err.Error())
	}
	application := app.NewApplication(context.TODO(), config)
	application.Setup()
}
