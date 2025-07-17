package main

import (
	"github.com/pasca-l/database-connection-manager/cli"
	"github.com/pasca-l/database-connection-manager/connection"
)

func main() {
	config := cli.NewConfig()
	cm := connection.NewConnectionManager()

	cli := cli.NewCli(cm, config)
	cli.HandleCommand()
}
