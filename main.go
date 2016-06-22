package main

import (
	"os"

	"github.com/intelsdi-x/snap/control/plugin"
	"github.com/raintank/snap-plugin-collector-ping/ping"
)

func main() {

	plugin.Start(
		ping.Meta(),
		ping.New(),
		os.Args[1],
	)
}
