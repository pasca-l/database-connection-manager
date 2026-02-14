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
	connectionManager connection.ConnectionManager
)

var rootCmd = &cobra.Command{
	Use:   "dbcm",
	Short: "Database Connection Manager",
	Long:  `A CLI tool to manage database connections and sessions for PostgreSQL and MySQL.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	cfg = config.NewConfig()
	connectionManager = connection.NewConnectionManager()
}

// loadConfig loads the configuration and connection manager state
// This should be called by commands that need the config (all except init)
func loadConfig() error {
	if _, err := os.Stat(cfg.Path); os.IsNotExist(err) {
		return fmt.Errorf("configuration not found\nRun 'dbcm init' to initialize")
	}

	state, err := cfg.Load()
	if err != nil {
		return fmt.Errorf("error loading configuration: %w", err)
	}

	if err := connectionManager.Load(state); err != nil {
		return fmt.Errorf("error loading connection manager state: %w", err)
	}

	return nil
}
