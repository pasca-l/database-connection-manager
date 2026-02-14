package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize configuration file",
	Long:  `Initialize the dbcm configuration file in the user's home directory.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := cfg.Init(); err != nil {
			return fmt.Errorf("error initializing configuration: %w", err)
		}
		fmt.Printf("Configuration initialized successfully at: %s\n", cfg.Path)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
