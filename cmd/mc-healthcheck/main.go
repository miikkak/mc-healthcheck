// Command mc-healthcheck checks whether a Minecraft server is responding, for
// use in container health checks. It supports Java Edition (Server List Ping
// over TCP) and Bedrock Edition (RakNet unconnected ping over UDP).
package main

import (
	"fmt"
	"os"

	"github.com/miikkak/mc-healthcheck/internal/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
