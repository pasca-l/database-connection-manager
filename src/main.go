package main

import (
	"fmt"
	"os"

	"github.com/pasca-l/database-connection-manager/cli"
	"github.com/pasca-l/database-connection-manager/connection"
)

func main() {
	config := cli.NewSaveConfig()

	cm := connection.NewConnectionManager(config.Path)
	if err := cm.Load(); err != nil {
		fmt.Printf("Error loading configuration: %v\n", err)
		os.Exit(1)
	}

	cli := cli.NewCli(*cm)
	cli.HandleCommand()
}
