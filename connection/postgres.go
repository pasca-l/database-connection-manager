package connection

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

type PostgreSQLConnector struct {
	Connection
}

func (p PostgreSQLConnector) TestConnection() bool {
	cmd := exec.Command("pg_isready", "-h", p.Host, "-p", strconv.Itoa(p.Port))
	err := cmd.Run()
	return err == nil
}

func (p PostgreSQLConnector) BuildCommand() *exec.Cmd {
	args := []string{
		"-h", p.Host,
		"-p", strconv.Itoa(p.Port),
		"-U", p.Username,
		"-d", p.Database,
	}

	cmd := exec.Command("psql", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Set PGPASSWORD environment variable if password is provided.
	if p.Password != "" {
		cmd.Env = append(os.Environ(), "PGPASSWORD="+p.Password)
	}

	return cmd
}

func (p PostgreSQLConnector) GetConnectionString() string {
	return fmt.Sprintf("postgresql://%s@%s:%d/%s", p.Username, p.Host, p.Port, p.Database)
}
