package cli

import (
	"flag"
	"fmt"

	"github.com/pasca-l/database-connection-manager/connection"
)

func ParseConnectionFlags(args []string) (*connection.Connection, error) {
	name := args[0]
	dbType := args[1]

	switch dbType {
	case "psql":
		return parsePostgreSQLFlags(name, args[2:])
	case "mysql":
		return parseMySQLFlags(name, args[2:])
	default:
		return nil, fmt.Errorf("unsupported connection type: %s", dbType)
	}
}

func parsePostgreSQLFlags(name string, flags []string) (*connection.Connection, error) {
	psqlFlags := flag.NewFlagSet("add psql", flag.ExitOnError)
	host := psqlFlags.String("host", "localhost", "Database host")
	hostShort := psqlFlags.String("h", "localhost", "Database host (shorthand)")
	port := psqlFlags.Int("port", 5432, "Database port")
	portShort := psqlFlags.Int("p", 5432, "Database port (shorthand)")
	database := psqlFlags.String("database", "", "Database name")
	databaseShort := psqlFlags.String("d", "", "Database name (shorthand)")
	username := psqlFlags.String("username", "", "Username")
	usernameShort := psqlFlags.String("U", "", "Username (shorthand)")
	password := psqlFlags.String("password", "", "Password")
	passwordShort := psqlFlags.String("w", "", "Password (shorthand)")

	psqlFlags.Parse(flags)

	// use shorthand values if main flags are empty
	if *host == "localhost" && *hostShort != "localhost" {
		*host = *hostShort
	}
	if *port == 5432 && *portShort != 5432 {
		*port = *portShort
	}
	if *database == "" && *databaseShort != "" {
		*database = *databaseShort
	}
	if *username == "" && *usernameShort != "" {
		*username = *usernameShort
	}
	if *password == "" && *passwordShort != "" {
		*password = *passwordShort
	}

	if *database == "" {
		return nil, fmt.Errorf("database name is required. Use -d or -database flag")
	}
	if *username == "" {
		return nil, fmt.Errorf("username is required. Use -U or -username flag")
	}

	return &connection.Connection{
		Name:     name,
		Type:     "psql",
		Host:     *host,
		Port:     *port,
		Database: *database,
		Username: *username,
		Password: *password,
	}, nil
}

func parseMySQLFlags(name string, flags []string) (*connection.Connection, error) {
	mysqlFlags := flag.NewFlagSet("add mysql", flag.ExitOnError)
	host := mysqlFlags.String("host", "localhost", "Database host")
	hostShort := mysqlFlags.String("h", "localhost", "Database host (shorthand)")
	port := mysqlFlags.Int("port", 3306, "Database port")
	portShort := mysqlFlags.Int("P", 3306, "Database port (shorthand)")
	database := mysqlFlags.String("database", "", "Database name")
	databaseShort := mysqlFlags.String("D", "", "Database name (shorthand)")
	username := mysqlFlags.String("username", "", "Username")
	usernameShort := mysqlFlags.String("u", "", "Username (shorthand)")
	password := mysqlFlags.String("password", "", "Password")
	passwordShort := mysqlFlags.String("p", "", "Password (shorthand)")

	mysqlFlags.Parse(flags)

	// use shorthand values if main flags are empty
	if *host == "localhost" && *hostShort != "localhost" {
		*host = *hostShort
	}
	if *port == 3306 && *portShort != 3306 {
		*port = *portShort
	}
	if *database == "" && *databaseShort != "" {
		*database = *databaseShort
	}
	if *username == "" && *usernameShort != "" {
		*username = *usernameShort
	}
	if *password == "" && *passwordShort != "" {
		*password = *passwordShort
	}

	if *database == "" {
		return nil, fmt.Errorf("database name is required. Use -D or -database flag")
	}
	if *username == "" {
		return nil, fmt.Errorf("username is required. Use -u or -username flag")
	}

	return &connection.Connection{
		Name:     name,
		Type:     "mysql",
		Host:     *host,
		Port:     *port,
		Database: *database,
		Username: *username,
		Password: *password,
	}, nil
}
