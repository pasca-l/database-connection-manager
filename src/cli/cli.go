package cli

import (
	"fmt"
	"os"

	"github.com/pasca-l/database-connection-manager/connection"
)

type ConnectionManagementCli interface {
	// function to handle commands
	HandleCommand()

	// function for executing commands
	handleList()
	handleAdd(args []string)
	handleRemove(args []string)
	handleConnect(args []string)
	handleSessions()
}

type Cli struct {
	ConnectionManager connection.ConnectionManager
	Config            Config
}

func NewCli(cm connection.ConnectionManager, config Config) ConnectionManagementCli {
	return &Cli{
		ConnectionManager: cm,
		Config:            config,
	}
}

func (c *Cli) HandleCommand() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]
	// check if config file exists, if not suggest running init
	if command != "init" {
		if _, err := os.Stat(c.Config.Path); os.IsNotExist(err) {
			fmt.Printf("Configuration not found\nRun 'dbcm init' to initialize\n")
			os.Exit(1)
		}
		state, err := c.Config.Load()
		if err != nil {
			fmt.Printf("Error loading configuration: %v\n", err)
			os.Exit(1)
		}
		if err := c.ConnectionManager.Load(state); err != nil {
			fmt.Printf("Error loading connection manager state: %v\n", err)
			os.Exit(1)
		}
	}

	switch command {
	case "init":
		c.handleInit()
	case "ls":
		c.handleList()
	case "add":
		c.handleAdd(os.Args[2:])
	case "remove":
		c.handleRemove(os.Args[2:])
	case "connect":
		c.handleConnect(os.Args[2:])
	case "sessions":
		c.handleSessions()
	default:
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Database Connection Manager (dbcm)")
	fmt.Println("Usage:")
	fmt.Println("  dbcm init                         Initialize configuration file")
	fmt.Println("  dbcm ls                           List all database connections and sessions")
	fmt.Println("  dbcm add <name> <type> [flags]    Add a new database connection")
	fmt.Println("  dbcm connect <name>               Connect to a database (creates/resumes session)")
	fmt.Println("  dbcm remove <name>                Remove a connection from management")
	fmt.Println("  dbcm sessions                     List active database sessions")
	fmt.Println()
	fmt.Println("PostgreSQL add flags:")
	fmt.Println("  -h, -host string       Database host (default: localhost)")
	fmt.Println("  -p, -port int          Database port (default: 5432)")
	fmt.Println("  -d, -database string   Database name (required)")
	fmt.Println("  -u, -username string   Username (required)")
	fmt.Println("  -w, -password string   Password (optional)")
	fmt.Println()
	fmt.Println("MySQL add flags:")
	fmt.Println("  -h, -host string       Database host (default: localhost)")
	fmt.Println("  -P, -port int          Database port (default: 3306)")
	fmt.Println("  -d, -database string   Database name (required)")
	fmt.Println("  -u, -username string   Username (required)")
	fmt.Println("  -p, -password string   Password (optional)")
}

func (c *Cli) handleInit() {
	if err := c.Config.Init(); err != nil {
		fmt.Printf("Error initializing configuration: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Configuration initialized successfully at: %s\n", c.Config.Path)
}

func (c *Cli) handleList() {
	fmt.Printf("%-15s %-8s %-20s %-8s %-15s %-12s %-8s\n", "NAME", "TYPE", "HOST", "PORT", "DATABASE", "STATUS", "SESSION")
	fmt.Println("---------------------------------------------------------------------------------------------")

	for _, conn := range c.ConnectionManager.Connections {
		status := "unreachable"
		if conn.TestConnection() {
			status = "reachable"
		}

		sessionStatus := ""
		sessions := c.ConnectionManager.GetSessionsByConnectionName(conn.Name)
		if len(sessions) > 0 {
			sessionStatus = "active*"
		}

		fmt.Printf("%-15s %-8s %-20s %-8d %-15s %-12s %-8s\n",
			conn.Name, conn.Type, conn.Host, conn.Port, conn.Database, status, sessionStatus)
	}
}

func (c *Cli) handleAdd(args []string) {
	if len(args) < 2 {
		fmt.Println("Error: connection name and type are required")
		fmt.Println("Usage: dbcm add <name> <type> [flags]")
		os.Exit(1)
	}

	conn, err := ParseConnectionFlags(args)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	if err := c.ConnectionManager.AddConnection(*conn, c.Config.Path); err != nil {
		fmt.Printf("Error adding connection: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("%s connection '%s' added successfully\n", conn.Type, conn.Name)
}

func (c *Cli) handleRemove(args []string) {
	if len(args) < 1 {
		fmt.Println("Error: connection name is required")
		fmt.Println("Usage: dbcm remove <name>")
		os.Exit(1)
	}

	name := args[0]

	if err := c.ConnectionManager.RemoveConnection(name, c.Config.Path); err != nil {
		fmt.Printf("Error removing connection: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Connection '%s' removed successfully\n", name)
}

func (c *Cli) handleConnect(args []string) {
	if len(args) < 1 {
		fmt.Println("Error: connection name is required")
		fmt.Println("Usage: dbcm connect <name>")
		os.Exit(1)
	}

	name := args[0]

	fmt.Printf("Connecting to '%s'\n", name)
	if err := c.ConnectionManager.Connect(name, c.Config.Path); err != nil {
		fmt.Printf("Error connecting to '%s': %v\n", name, err)
		os.Exit(1)
	}
}

func (c *Cli) handleSessions() {
	fmt.Printf("%-15s %-8s %-20s\n", "CONNECTION", "PID", "STARTED")
	fmt.Println("-------------------------------------------")

	for _, session := range c.ConnectionManager.Sessions {
		fmt.Printf("%-15s %-8d %-20s\n",
			session.ConnectionName, session.PID, session.Started.Format("2006-01-02 15:04"))
	}
}
