package cmd

import (
	"fmt"

	"github.com/pasca-l/database-connection-manager/connection"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add <name> <type>",
	Short: "Add a new database connection",
	Long: `Add a new database connection to the manager.
Supported types: psql (PostgreSQL), mysql (MySQL)`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := loadConfig(); err != nil {
			return err
		}

		name := args[0]
		dbType := args[1]

		var conn *connection.Connection
		var err error

		switch dbType {
		case "psql":
			conn, err = createPostgreSQLConnection(cmd, name)
		case "mysql":
			conn, err = createMySQLConnection(cmd, name)
		default:
			return fmt.Errorf("unsupported connection type: %s", dbType)
		}

		if err != nil {
			return err
		}

		if err := connectionManager.AddConnection(*conn, cfg.Path); err != nil {
			return fmt.Errorf("error adding connection: %w", err)
		}

		fmt.Printf("%s connection '%s' added successfully\n", conn.Type, conn.Name)
		return nil
	},
}

func createPostgreSQLConnection(cmd *cobra.Command, name string) (*connection.Connection, error) {
	host, _ := cmd.Flags().GetString("host")
	port, _ := cmd.Flags().GetInt("port")
	database, _ := cmd.Flags().GetString("database")
	username, _ := cmd.Flags().GetString("username")
	password, _ := cmd.Flags().GetString("password")

	// Use default port for PostgreSQL if not explicitly set
	if !cmd.Flags().Changed("port") && !cmd.Flags().Changed("p") {
		port = 5432
	}

	if database == "" {
		return nil, fmt.Errorf("database name is required. Use -d or --database flag")
	}
	if username == "" {
		return nil, fmt.Errorf("username is required. Use -U or --username flag")
	}

	return &connection.Connection{
		Name:     name,
		Type:     "psql",
		Host:     host,
		Port:     port,
		Database: database,
		Username: username,
		Password: password,
	}, nil
}

func createMySQLConnection(cmd *cobra.Command, name string) (*connection.Connection, error) {
	host, _ := cmd.Flags().GetString("host")

	// For MySQL, check both -P (uppercase) and -p flags
	port := 3306
	if cmd.Flags().Changed("port") || cmd.Flags().Changed("p") {
		port, _ = cmd.Flags().GetInt("port")
	}

	// For MySQL, check both -D and -d flags for database
	database, _ := cmd.Flags().GetString("database")

	// For MySQL, check both -u and -U flags for username
	username, _ := cmd.Flags().GetString("username")

	password, _ := cmd.Flags().GetString("password")

	if database == "" {
		return nil, fmt.Errorf("database name is required. Use -D or --database flag")
	}
	if username == "" {
		return nil, fmt.Errorf("username is required. Use -u or --username flag")
	}

	return &connection.Connection{
		Name:     name,
		Type:     "mysql",
		Host:     host,
		Port:     port,
		Database: database,
		Username: username,
		Password: password,
	}, nil
}

func init() {
	rootCmd.AddCommand(addCmd)

	// Common flags that work for both PostgreSQL and MySQL
	// Note: Using long-form --host instead of -h to avoid conflict with Cobra's help flag
	addCmd.Flags().String("host", "localhost", "Database host")
	addCmd.Flags().IntP("port", "p", 0, "Database port (default: psql=5432, mysql=3306)")
	addCmd.Flags().StringP("database", "d", "", "Database name (required)")
	addCmd.Flags().StringP("username", "U", "", "Username (required, -U for psql, -u for mysql)")
	addCmd.Flags().StringP("password", "w", "", "Password (optional)")
}
