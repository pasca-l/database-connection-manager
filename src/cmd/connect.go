package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var connectCmd = &cobra.Command{
	Use:   "connect <name>",
	Short: "Connect to a database (creates/resumes session)",
	Long:  `Connect to a database using the specified connection name. Creates a new session or resumes an existing one.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := loadConfig(); err != nil {
			return err
		}

		name := args[0]

		fmt.Printf("Connecting to '%s'\n", name)
		if err := connectionManager.Connect(name, cfg.Path); err != nil {
			return fmt.Errorf("error connecting to '%s': %w", name, err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(connectCmd)
}
