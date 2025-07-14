package connection

import (
	"encoding/json"
	"os"
	"time"
)

type ConnectionManager struct {
	Connections Connections `json:"connections,omitempty"`
	Sessions    Sessions    `json:"sessions,omitempty"`
	savePath    string      `json:"-"`
}

func NewConnectionManager(savePath string) *ConnectionManager {
	return &ConnectionManager{
		savePath: savePath,
	}
}

func (cm *ConnectionManager) Load() error {
	// load connections and sessions
	if _, err := os.Stat(cm.savePath); !os.IsNotExist(err) {
		data, err := os.ReadFile(cm.savePath)
		if err != nil {
			return err
		}
		if err := json.Unmarshal(data, &cm); err != nil {
			return err
		}
	}
	return nil
}

func (cm *ConnectionManager) Save() error {
	data, err := json.MarshalIndent(cm, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(cm.savePath, data, 0644)
}

func (cm *ConnectionManager) Clear() error {
	cm.Connections = make(Connections, 0)
	cm.Sessions = make(Sessions, 0)
	return cm.Save()
}

func (cm *ConnectionManager) AddConnection(conn Connection) error {
	if err := cm.Connections.AddConnection(conn); err != nil {
		return err
	}
	return cm.Save()
}

func (cm *ConnectionManager) RemoveConnection(name string) error {
	if err := cm.Connections.RemoveConnection(name); err != nil {
		return err
	}
	if err := cm.Sessions.RemoveSessionsByConnectionName(name); err != nil {
		return err
	}
	return cm.Save()
}

func (cm *ConnectionManager) GetConnection(name string) (*Connection, error) {
	return cm.Connections.GetConnection(name)
}

func (cm *ConnectionManager) CleanupSessions() error {
	cm.Sessions.CleanupSessions()
	return cm.Save()
}

func (cm *ConnectionManager) GetSessionsByConnectionName(name string) Sessions {
	return cm.Sessions.GetSessionsByConnectionName(name)
}

func (cm *ConnectionManager) AddSession(connectionName string, pid int) error {
	newSession := Session{
		ConnectionName: connectionName,
		PID:            pid,
		Started:        time.Now(),
		Active:         true,
	}
	cm.Sessions.AddSession(newSession)

	return cm.Save()
}
