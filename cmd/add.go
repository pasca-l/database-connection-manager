package cmd

import (
	"errors"
	"fmt"

	"github.com/pasca-l/database-connection-manager/connection"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new database connection",
	Long:  `Add a new database connection to the manager. Use subcommands for specific database types.`,
}

var (
	ErrDatabaseNameRequired = errors.New("database name is required (use -d or --database)")
	ErrUsernameRequired     = errors.New("username is required (use -U or --username)")
)

var addPsqlCmd = &cobra.Command{
	Use:   "psql <name>",
	Short: "Add a PostgreSQL connection",
	Long:  `Add a PostgreSQL database connection with PostgreSQL-specific flags.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		flags, err := parseFlags(cmd)
		if err != nil {
			return err
		}

		// Use default port for PostgreSQL if not explicitly set.
		if !cmd.Flags().Changed("port") {
			flags.Port = 5432
		}

		if flags.Database == "" {
			return ErrDatabaseNameRequired
		}
		if flags.Username == "" {
			return ErrUsernameRequired
		}

		conn := &connection.Connection{
			Name:     name,
			Type:     "psql",
			Host:     flags.Host,
			Port:     flags.Port,
			Database: flags.Database,
			Username: flags.Username,
			Password: flags.Password,
		}

		if err := connectionManager.AddConnection(*conn); err != nil {
			return fmt.Errorf("error adding connection: %w", err)
		}

		fmt.Printf("PostgreSQL connection '%s' added successfully\n", conn.Name)
		return nil
	},
}

var addMysqlCmd = &cobra.Command{
	Use:   "mysql <name>",
	Short: "Add a MySQL connection",
	Long:  `Add a MySQL database connection with MySQL-specific flags.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		flags, err := parseFlags(cmd)
		if err != nil {
			return err
		}

		// Use default port for MySQL if not explicitly set.
		if !cmd.Flags().Changed("port") {
			flags.Port = 3306
		}

		if flags.Database == "" {
			return ErrDatabaseNameRequired
		}
		if flags.Username == "" {
			return ErrUsernameRequired
		}

		conn := &connection.Connection{
			Name:     name,
			Type:     "mysql",
			Host:     flags.Host,
			Port:     flags.Port,
			Database: flags.Database,
			Username: flags.Username,
			Password: flags.Password,
		}

		if err := connectionManager.AddConnection(*conn); err != nil {
			return fmt.Errorf("error adding connection: %w", err)
		}

		fmt.Printf("MySQL connection '%s' added successfully\n", conn.Name)
		return nil
	},
}

type dnsFlags struct {
	Host     string
	Port     int
	Database string
	Username string
	Password string
}

func parseFlags(cmd *cobra.Command) (*dnsFlags, error) {
	host, err := cmd.Flags().GetString("host")
	if err != nil {
		return nil, fmt.Errorf("error reading host flag: %w", err)
	}
	port, err := cmd.Flags().GetInt("port")
	if err != nil {
		return nil, fmt.Errorf("error reading port flag: %w", err)
	}
	database, err := cmd.Flags().GetString("database")
	if err != nil {
		return nil, fmt.Errorf("error reading database flag: %w", err)
	}
	username, err := cmd.Flags().GetString("username")
	if err != nil {
		return nil, fmt.Errorf("error reading username flag: %w", err)
	}
	password, err := cmd.Flags().GetString("password")
	if err != nil {
		return nil, fmt.Errorf("error reading password flag: %w", err)
	}

	return &dnsFlags{
		Host:     host,
		Port:     port,
		Database: database,
		Username: username,
		Password: password,
	}, nil
}

func init() {
	rootCmd.AddCommand(addCmd)

	// Add database-specific subcommands.
	addCmd.AddCommand(addPsqlCmd)
	addCmd.AddCommand(addMysqlCmd)

	// Disable automatic help flag to allow -h for host.
	addPsqlCmd.Flags().BoolP("help", "", false, "Help for psql")
	addMysqlCmd.Flags().BoolP("help", "", false, "Help for mysql")

	// PostgreSQL-specific flags (matching psql conventions).
	addPsqlCmd.Flags().StringP("host", "h", "localhost", "Database host")
	addPsqlCmd.Flags().IntP("port", "p", 5432, "Database port")
	addPsqlCmd.Flags().StringP("database", "d", "", "Database name (required)")
	addPsqlCmd.Flags().StringP("username", "U", "", "Username (required)")
	addPsqlCmd.Flags().StringP("password", "w", "", "Password (optional)")

	// MySQL-specific flags (matching mysql conventions).
	addMysqlCmd.Flags().StringP("host", "h", "localhost", "Database host")
	addMysqlCmd.Flags().IntP("port", "P", 3306, "Database port")
	addMysqlCmd.Flags().StringP("database", "D", "", "Database name (required)")
	addMysqlCmd.Flags().StringP("username", "u", "", "Username (required)")
	addMysqlCmd.Flags().StringP("password", "p", "", "Password (optional)")
}
