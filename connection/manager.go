package connection

import (
	"encoding/json"
)

type ConfigStore interface {
	Load() ([]byte, error)
	Save(data []byte) error
}

type ConnectionManager struct {
	Connections Connections `json:"connections,omitempty"`
	config      ConfigStore
}

func NewConnectionManager(config ConfigStore) *ConnectionManager {
	return &ConnectionManager{
		Connections: make(Connections, 0),
		config:      config,
	}
}

func (cm *ConnectionManager) Load() error {
	data, err := cm.config.Load()
	if err != nil {
		return err
	}
	return json.Unmarshal(data, cm)
}

func (cm *ConnectionManager) Save() error {
	data, err := json.MarshalIndent(cm, "", "  ")
	if err != nil {
		return err
	}
	return cm.config.Save(data)
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
	return cm.Save()
}

func (cm *ConnectionManager) GetConnection(name string) (Connection, error) {
	return cm.Connections.GetConnection(name)
}

func (cm *ConnectionManager) GetAllConnections() Connections {
	return cm.Connections
}
