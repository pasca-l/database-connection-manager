package connection

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"
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
	conn, err := cm.Connections.GetConnection(name)
	if err != nil {
		return err
	}
	cmd, err := conn.ConnectCmd()
	if err != nil {
		return fmt.Errorf("failed to connect to '%s': %v", name, err)
	}

	// add session to manager after starting connection command
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start command: %v", err)
	}
	cm.Sessions.AddSession(NewSession(name, cmd.Process.Pid))
	if err := cm.Save(configPath); err != nil {
		return err
	}

	// handle termination signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		// forward signal to the command process
		_ = cmd.Process.Signal(syscall.SIGINT)
	}()

	// remove session from manager after connection command finishes
	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("failed to wait for command: %v", err)
	}
	cm.Sessions.RemoveSession(name, cmd.Process.Pid)
	if err := cm.Save(configPath); err != nil {
		return err
	}
	return nil
}
