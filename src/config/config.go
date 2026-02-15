package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

var (
	ErrConfigAlreadyExists = errors.New("configuration already exists")
	ErrConfigNotFound      = errors.New("configuration not found\nRun 'dbcm init' to initialize")
)

type Config struct {
	Path string
}

func NewConfig() Config {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Error getting home directory: %v\n", err)
		os.Exit(1)
	}

	configDir := filepath.Join(homeDir, ".dbcm")
	configPath := filepath.Join(configDir, "state.json")

	return Config{
		Path: configPath,
	}
}

func (c Config) Init() error {
	if _, err := os.Stat(c.Path); err == nil {
		return ErrConfigAlreadyExists
	}

	// create config directory if necessary
	if err := os.MkdirAll(filepath.Dir(c.Path), 0750); err != nil {
		return fmt.Errorf("error creating config directory: %w", err)
	}
	// create an empty config file
	if err := os.WriteFile(c.Path, []byte("{}"), 0600); err != nil {
		return fmt.Errorf("error creating config file: %w", err)
	}

	return nil
}

func (c Config) Load() ([]byte, error) {
	if _, err := os.Stat(c.Path); os.IsNotExist(err) {
		return nil, ErrConfigNotFound
	}
	return os.ReadFile(c.Path)
}

func (c Config) Save(data []byte) error {
	return os.WriteFile(c.Path, data, 0600)
}
