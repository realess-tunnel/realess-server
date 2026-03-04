package main

import (
	"realess-server/internal/commands"
)

var version = "dev"

func main() {
	commands.CreateRootCmd(version).Execute()
}
