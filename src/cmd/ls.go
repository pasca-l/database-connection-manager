package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List all database connections and sessions",
	Long:  `Display a table of all database connections with their status and active sessions.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := loadConfig(); err != nil {
			return err
		}

		fmt.Printf("%-15s %-8s %-20s %-8s %-15s %-12s %-8s\n", "NAME", "TYPE", "HOST", "PORT", "DATABASE", "STATUS", "SESSION")
		fmt.Println("---------------------------------------------------------------------------------------------")

		for _, conn := range connectionManager.Connections {
			status := "unreachable"
			if conn.TestConnection() {
				status = "reachable"
			}

			sessionStatus := ""
			sessions := connectionManager.GetSessionsByConnectionName(conn.Name)
			if len(sessions) > 0 {
				sessionStatus = "active*"
			}

			fmt.Printf("%-15s %-8s %-20s %-8d %-15s %-12s %-8s\n",
				conn.Name, conn.Type, conn.Host, conn.Port, conn.Database, status, sessionStatus)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(lsCmd)
}
