package cmd

import (
	"fmt"
	"os"

	"github.com/pasca-l/database-connection-manager/config"
	"github.com/pasca-l/database-connection-manager/connection"
	"github.com/spf13/cobra"
)

var (
	cfg               config.Config
	connectionManager *connection.ConnectionManager
)

var rootCmd = &cobra.Command{
	Use:   "dbcm",
	Short: "Database Connection Manager",
	Long:  `A CLI tool to manage database connections for PostgreSQL and MySQL.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Skip loading config for init command.
		if cmd.Name() == "init" {
			return nil
		}
		// Load configuration for all other commands.
		return connectionManager.Load()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	cfg = config.NewConfig()
	connectionManager = connection.NewConnectionManager(cfg)
}
