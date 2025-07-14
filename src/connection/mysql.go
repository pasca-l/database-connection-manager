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
	}

	if m.Database != "" {
		args = append(args, m.Database)
	}
	if m.Password != "" {
		args = append(args, fmt.Sprintf("-p%s", m.Password))
	}

	cmd := exec.Command("mysql", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd
}

func (m MySQLConnector) GetConnectionString() string {
	return fmt.Sprintf("mysql://%s@%s:%d/%s", m.Username, m.Host, m.Port, m.Database)
}
