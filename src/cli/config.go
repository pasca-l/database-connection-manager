package cli

import (
	"fmt"
	"os"
	"path/filepath"
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
		return fmt.Errorf("configuration already exists")
	}

	// create config directory if necessary
	if err := os.MkdirAll(filepath.Dir(c.Path), 0755); err != nil {
		return fmt.Errorf("error creating config directory: %v", err)
	}
	// create an empty config file
	if err := os.WriteFile(c.Path, []byte("{}"), 0644); err != nil {
		return fmt.Errorf("error creating config file: %v", err)
	}

	return nil
}

func (c Config) Load() ([]byte, error) {
	return os.ReadFile(c.Path)
}

func (c Config) Save(data []byte) error {
	return os.WriteFile(c.Path, data, 0644)
}
