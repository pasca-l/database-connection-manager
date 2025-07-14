package connection

import (
	"fmt"
	"os/exec"
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

func (cs *Connections) GetConnection(name string) (*Connection, error) {
	for _, conn := range *cs {
		if conn.Name == name {
			return &conn, nil
		}
	}
	return nil, fmt.Errorf("connection not found: %s", name)
}

func (cs *Connections) AddConnection(conn Connection) error {
	for _, existing := range *cs {
		if existing.Name == conn.Name {
			return fmt.Errorf("connection already exists: %s", conn.Name)
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
	return fmt.Errorf("connection not found: %s", name)
}
