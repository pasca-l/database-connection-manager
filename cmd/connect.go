package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var connectCmd = &cobra.Command{
	Use:   "connect <name>",
	Short: "Connect to a database",
	Long:  `Connect to a database using the specified connection name and launch the native database client.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		conn, err := connectionManager.GetConnection(name)
		if err != nil {
			return fmt.Errorf("connection '%s' not found: %w", name, err)
		}
		dbCmd, err := conn.ConnectCmd()
		if err != nil {
			return fmt.Errorf("failed to build connection command: %w", err)
		}

		fmt.Printf("Connecting to '%s'...\n", name)
		if err := dbCmd.Run(); err != nil {
			return fmt.Errorf("error connecting to '%s': %w", name, err)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(connectCmd)
}
