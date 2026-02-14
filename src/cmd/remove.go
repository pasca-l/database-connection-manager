package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:   "remove <name>",
	Short: "Remove a connection from management",
	Long:  `Remove a database connection from the connection manager.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := loadConfig(); err != nil {
			return err
		}

		name := args[0]

		if err := connectionManager.RemoveConnection(name, cfg.Path); err != nil {
			return fmt.Errorf("error removing connection: %w", err)
		}

		fmt.Printf("Connection '%s' removed successfully\n", name)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)
}
