package connection

import (
	"fmt"
	"os"
	"slices"
	"syscall"
	"time"
)

type Sessions []Session
type Session struct {
	ConnectionName string    `json:"connection_name"`
	PID            int       `json:"pid"`
	Started        time.Time `json:"started"`
	Active         bool      `json:"active"`
}

func (ss *Sessions) GetSessionsByConnectionName(name string) Sessions {
	var filteredSessions Sessions
	for _, session := range *ss {
		if session.ConnectionName == name && session.IsProcessAlive() {
			filteredSessions = append(filteredSessions, session)
		}
	}
	return filteredSessions
}

func (ss *Sessions) AddSession(s Session) {
	*ss = append(*ss, s)
}

func (ss *Sessions) RemoveSessionsByConnectionName(name string) error {
	for i := len(*ss) - 1; i >= 0; i-- {
		session := (*ss)[i]
		if session.ConnectionName == name {
			// kill the process if it is still alive
			if session.IsProcessAlive() {
				err := session.Kill()
				if err != nil {
					return fmt.Errorf("failed to kill session of connection %s: %v", name, err)
				}
			}
			*ss = slices.Delete(*ss, i, i+1)
		}
	}
	return nil
}

func (ss *Sessions) CleanupSessions() {
	var activeSessions Sessions
	for _, session := range *ss {
		if session.IsProcessAlive() {
			activeSessions = append(activeSessions, session)
		}
	}
	*ss = activeSessions
}

func (ss *Sessions) AnyActive() bool {
	return slices.ContainsFunc(*ss, func(s Session) bool { return s.Active })
}

func NewSession(connectionName string, pid int) Session {
	return Session{
		ConnectionName: connectionName,
		PID:            pid,
		Started:        time.Now(),
		Active:         true,
	}
}

func (s Session) IsProcessAlive() bool {
	if s.PID <= 0 {
		return false
	}
	process, err := os.FindProcess(s.PID)
	if err != nil {
		return false
	}
	// send signal 0 to check if process exists
	err = process.Signal(syscall.Signal(0))
	return err == nil
}

func (s Session) Kill() error {
	process, err := os.FindProcess(s.PID)
	if err != nil {
		return err
	}
	return process.Kill()
}

func (s Session) Continue() error {
	process, err := os.FindProcess(s.PID)
	if err != nil {
		return err
	}
	return process.Signal(syscall.SIGCONT)
}
