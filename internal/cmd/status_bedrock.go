package cmd

import (
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/sandertv/go-raknet"
	"github.com/spf13/cobra"
)

var statusBedrockCmd = &cobra.Command{
	Use:   "status-bedrock",
	Short: "Check a Bedrock Edition server via a RakNet unconnected ping (UDP)",
	Example: `
mc-healthcheck status-bedrock --host localhost --port 19132 --timeout 5s
`,
	RunE: func(cmd *cobra.Command, _ []string) error {
		host, _ := cmd.Flags().GetString("host")
		port, _ := cmd.Flags().GetInt("port")
		timeout, _ := cmd.Flags().GetDuration("timeout")

		if timeout <= 0 {
			return fmt.Errorf("timeout must be positive, got %s", timeout)
		}

		addr := net.JoinHostPort(host, strconv.Itoa(port))
		response, err := raknet.PingTimeout(addr, timeout)
		if err != nil {
			return fmt.Errorf("status-bedrock check failed: %w", err)
		}
		if len(response) == 0 {
			return fmt.Errorf("status-bedrock check failed: empty response from %s", addr)
		}
		return nil
	},
}

func init() {
	statusBedrockCmd.Flags().String("host", "localhost", "server hostname")
	statusBedrockCmd.Flags().Int("port", 19132, "server port")
	statusBedrockCmd.Flags().Duration("timeout", 5*time.Second, "read timeout")
	RootCmd.AddCommand(statusBedrockCmd)
}
