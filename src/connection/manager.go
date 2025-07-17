package connection

import (
	"encoding/json"
	"fmt"
	"os"
)

type ConnectionManager struct {
	Connections Connections `json:"connections,omitempty"`
	Sessions    Sessions    `json:"sessions,omitempty"`
}

func NewConnectionManager() ConnectionManager {
	return ConnectionManager{}
}

func (cm *ConnectionManager) Load(data []byte) error {
	return json.Unmarshal(data, &cm)
}

func (cm *ConnectionManager) Save(configPath string) error {
	data, err := json.MarshalIndent(cm, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(configPath, data, 0644)
}

func (cm *ConnectionManager) Clear(configPath string) error {
	cm.Connections = make(Connections, 0)
	cm.Sessions = make(Sessions, 0)
	return cm.Save(configPath)
}

func (cm *ConnectionManager) AddConnection(conn Connection, configPath string) error {
	if err := cm.Connections.AddConnection(conn); err != nil {
		return err
	}
	return cm.Save(configPath)
}

func (cm *ConnectionManager) RemoveConnection(name string, configPath string) error {
	if err := cm.Connections.RemoveConnection(name); err != nil {
		return err
	}
	if err := cm.Sessions.RemoveSessionsByConnectionName(name); err != nil {
		return err
	}
	return cm.Save(configPath)
}

func (cm *ConnectionManager) GetSessionsByConnectionName(name string) Sessions {
	return cm.Sessions.GetSessionsByConnectionName(name)
}

func (cm *ConnectionManager) Connect(name string, configPath string) error {
	// if any non-active sessions exist, resume the first one found
	sessions := cm.Sessions.GetSessionsByConnectionName(name)
	for _, session := range sessions {
		if session.Active {
			continue
		}
		// bring the existing process to foreground
		session.Continue()
		return cm.Save(configPath)
	}

	// if no non-active sessions found, start a new connection
	conn, err := cm.Connections.GetConnection(name)
	if err != nil {
		return err
	}
	pid, err := conn.Connect()
	if err != nil {
		return fmt.Errorf("failed to connect to '%s': %v", name, err)
	}

	// add session to tracking
	cm.Sessions.AddSession(NewSession(name, pid))
	return cm.Save(configPath)
}
