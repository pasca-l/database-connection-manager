package connection

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

type MySQLConnector struct {
	Connection
}

func (m MySQLConnector) TestConnection() bool {
	cmd := exec.Command("mysqladmin", "-h", m.Host, "-P", strconv.Itoa(m.Port), "-u", m.Username, "-p"+m.Password, "ping")
	err := cmd.Run()
	return err == nil
}

func (m MySQLConnector) BuildCommand() *exec.Cmd {
	args := []string{
		"-h", m.Host,
		"-P", strconv.Itoa(m.Port),
		"-u", m.Username,
		m.Database,
	}

	cmd := exec.Command("mysql", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Set MYSQL_PWD environment variable if password is provided.
	if m.Password != "" {
		cmd.Env = append(os.Environ(), fmt.Sprintf("MYSQL_PWD=%s", m.Password))
	}

	return cmd
}

func (m MySQLConnector) GetConnectionString() string {
	return fmt.Sprintf("mysql://%s@%s:%d/%s", m.Username, m.Host, m.Port, m.Database)
}
