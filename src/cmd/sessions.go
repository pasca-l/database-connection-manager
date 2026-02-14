package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var sessionsCmd = &cobra.Command{
	Use:   "sessions",
	Short: "List active database sessions",
	Long:  `Display a table of all active database sessions with their PIDs and start times.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := loadConfig(); err != nil {
			return err
		}

		fmt.Printf("%-15s %-8s %-20s\n", "CONNECTION", "PID", "STARTED")
		fmt.Println("-------------------------------------------")

		for _, session := range connectionManager.Sessions {
			fmt.Printf("%-15s %-8d %-20s\n",
				session.ConnectionName, session.PID, session.Started.Format("2006-01-02 15:04"))
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(sessionsCmd)
}
