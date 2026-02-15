package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List all database connections",
	Long:  `Display a table of all database connections with their status.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("%-15s %-8s %-20s %-8s %-15s %-12s\n", "NAME", "TYPE", "HOST", "PORT", "DATABASE", "STATUS")
		fmt.Println("--------------------------------------------------------------------------------")

		connections := connectionManager.GetAllConnections()
		for _, conn := range connections {
			status := "unreachable"
			if conn.TestConnection() {
				status = "reachable"
			}

			fmt.Printf("%-15s %-8s %-20s %-8d %-15s %-12s\n",
				conn.Name, conn.Type, conn.Host, conn.Port, conn.Database, status)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(lsCmd)
}
