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
}

func NewCli(cm connection.ConnectionManager) ConnectionManagementCli {
	return &Cli{
		ConnectionManager: cm,
	}
}

func (c *Cli) HandleCommand() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
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
	case "-h", "--help", "help":
		printUsage()
	default:
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Database Connection Manager (dbcm)")
	fmt.Println("Usage:")
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

func (c *Cli) handleList() {
	if len(c.ConnectionManager.Connections) == 0 {
		fmt.Println("No connections configured")
		return
	}

	fmt.Printf("%-15s %-8s %-20s %-8s %-15s %-12s %-8s\n", "NAME", "TYPE", "HOST", "PORT", "DATABASE", "STATUS", "SESSION")
	fmt.Println("----------------------------------------------------------------------------------------")

	for _, conn := range c.ConnectionManager.Connections {
		status := "unreachable"
		connector := conn.GetConnector()
		if connector != nil && connector.TestConnection() {
			status = "reachable"
		}

		sessionStatus := ""
		sessions := c.ConnectionManager.GetSessionsByConnectionName(conn.Name)
		if len(sessions) > 0 {
			sessionStatus = "active"
			if sessions.AnyActive() {
				sessionStatus = "active*"
			}
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

	name := args[0]
	connType := args[1]
	conn, err := ParseConnectionFlags(connType, args[2:])
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	conn.Name = name
	if err := c.ConnectionManager.AddConnection(*conn); err != nil {
		fmt.Printf("Error adding connection: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("%s connection '%s' added successfully\n", connType, name)
}

func (c *Cli) handleRemove(args []string) {
	if len(args) < 1 {
		fmt.Println("Error: connection name is required")
		fmt.Println("Usage: dbcm remove <name>")
		os.Exit(1)
	}

	name := args[0]

	if err := c.ConnectionManager.RemoveConnection(name); err != nil {
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

	// check for existing non-active sessions
	// if any, resume the first one found
	sessions := c.ConnectionManager.GetSessionsByConnectionName(name)
	for _, session := range sessions {
		if session.Active {
			continue
		}
		// bring the existing process to foreground
		session.Continue()

		fmt.Printf("Resuming existing session for '%s' (PID: %d)\n", name, session.PID)
		err := c.ConnectionManager.Save()
		if err != nil {
			fmt.Printf("Error saving session: %v\n", err)
			os.Exit(1)
		}
		return
	}

	// find the specified connection
	conn, err := c.ConnectionManager.GetConnection(name)
	if err != nil {
		fmt.Printf("Connection '%s' not found: %v\n", name, err)
		os.Exit(1)
	}
	fmt.Printf("Connecting to %s database '%s'...\n", conn.Type, name)

	// connect using the specified connector
	connector := conn.GetConnector()
	if connector == nil {
		fmt.Printf("Error: unsupported connection type '%s'\n", conn.Type)
		os.Exit(1)
	}
	cmd := connector.BuildCommand()
	if cmd == nil {
		fmt.Printf("Error: unable to build command for connection type '%s'\n", conn.Type)
		os.Exit(1)
	}
	if err := cmd.Start(); err != nil {
		fmt.Printf("Error starting %s: %v\n", conn.Type, err)
		os.Exit(1)
	}
	// wait for the command to complete
	if err := cmd.Wait(); err != nil {
		fmt.Printf("%s exited with error: %v\n", conn.Type, err)
	}

	// add session to tracking
	c.ConnectionManager.AddSession(name, cmd.Process.Pid)
	fmt.Printf("Started new session for '%s' (PID: %d)\n", name, cmd.Process.Pid)
}

func (c *Cli) handleSessions() {
	fmt.Println("sessions: ", c.ConnectionManager.Sessions)
	if len(c.ConnectionManager.Sessions) == 0 {
		fmt.Println("No active sessions")
		return
	}

	fmt.Printf("%-15s %-8s %-20s %-8s\n", "CONNECTION", "PID", "STARTED", "ACTIVE")
	fmt.Println("-------------------------------------------------------")

	for _, session := range c.ConnectionManager.Sessions {
		active := ""
		if session.Active {
			active = "*"
		}

		fmt.Printf("%-15s %-8d %-20s %-8s\n",
			session.ConnectionName, session.PID, session.Started.Format("2006-01-02 15:04"), active)
	}
}
