package cli

import (
	"fmt"
	"os"
	"path/filepath"
)

type SaveConfig struct {
	Path string
}

func NewSaveConfig() SaveConfig {
	workingDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting working directory: %v\n", err)
		os.Exit(1)
	}

	configDir := filepath.Join(workingDir, ".dbcm")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		fmt.Printf("Error creating config directory: %v\n", err)
		os.Exit(1)
	}

	configPath := filepath.Join(configDir, "state.json")
	return SaveConfig{
		Path: configPath,
	}
}
