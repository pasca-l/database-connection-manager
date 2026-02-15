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
	if err := json.Unmarshal(data, cm); err != nil {
		return err
	}
	return cm.Connections.DecryptPasswords()
}

func (cm *ConnectionManager) Save() error {
	// Create a copy of connections to encrypt passwords.
	encryptedConnections := make(Connections, len(cm.Connections))
	copy(encryptedConnections, cm.Connections)
	if err := encryptedConnections.EncryptPasswords(); err != nil {
		return err
	}

	data, err := json.MarshalIndent(
		&ConnectionManager{Connections: encryptedConnections}, "", "  ",
	)
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
