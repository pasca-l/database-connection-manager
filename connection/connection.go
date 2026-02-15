package connection

import (
	"errors"
	"fmt"
	"os/exec"
)

var (
	ErrUnsupportedConnectionType   = errors.New("unsupported connection type")
	ErrConnectionAlreadyRegistered = errors.New("connection already registered")
	ErrConnectionNotFound          = errors.New("connection not found")
	ErrBuildCommandFailed          = errors.New("failed to build command for connection")
)

type Connector interface {
	TestConnection() bool
	BuildCommand() *exec.Cmd
	GetConnectionString() string
}

type Connections []Connection
type Connection struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Database string `json:"database"`
	Username string `json:"username"`
	Password string `json:"password,omitempty"`
}

func (c Connection) GetConnector() Connector {
	switch c.Type {
	case "psql":
		return PostgreSQLConnector{c}
	case "mysql":
		return MySQLConnector{c}
	default:
		return nil
	}
}

func (c Connection) TestConnection() bool {
	connector := c.GetConnector()
	if connector == nil {
		return false
	}
	return connector.TestConnection()
}

func (c Connection) ConnectCmd() (*exec.Cmd, error) {
	connector := c.GetConnector()
	if connector == nil {
		return nil, fmt.Errorf("%w: %s", ErrUnsupportedConnectionType, c.Type)
	}
	cmd := connector.BuildCommand()
	if cmd == nil {
		return nil, fmt.Errorf("%w: %s", ErrBuildCommandFailed, c.Name)
	}
	return cmd, nil
}

func (cs *Connections) GetConnection(name string) (Connection, error) {
	for _, conn := range *cs {
		if conn.Name == name {
			return conn, nil
		}
	}
	return Connection{}, fmt.Errorf("%w: %s", ErrConnectionNotFound, name)
}

func (cs *Connections) AddConnection(conn Connection) error {
	for _, existing := range *cs {
		if existing.Name == conn.Name {
			return fmt.Errorf("%w: %s", ErrConnectionAlreadyRegistered, conn.Name)
		}
	}
	*cs = append(*cs, conn)
	return nil
}

func (cs *Connections) RemoveConnection(name string) error {
	for i, conn := range *cs {
		if conn.Name == name {
			*cs = append((*cs)[:i], (*cs)[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("%w: %s", ErrConnectionNotFound, name)
}
