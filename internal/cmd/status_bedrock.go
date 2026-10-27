package cmd

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/miikkak/mc-healthcheck/internal/nethernet"
	"github.com/sandertv/go-raknet"
	"github.com/spf13/cobra"
)

const (
	transportAuto      = "auto"
	transportNetherNet = "nethernet"
	transportRakNet    = "raknet"
)

var statusBedrockCmd = &cobra.Command{
	Use:   "status-bedrock",
	Short: "Check a Bedrock Edition server via NetherNet signaling (TCP) and/or a RakNet ping (UDP)",
	Long: `Checks a Bedrock listener (e.g. Geyser). Bedrock 26.60+ clients are
NetherNet-only, so the default transport "auto" checks NetherNet signaling
(GET /v1/join on the TCP port) first and falls back to a RakNet unconnected
ping (UDP). Use --transport to pin one. The NetherNet check confirms the
signaling endpoint only, not the UDP gameplay leg. Signaling is plain HTTP
unless --tls is set.`,
	Example: `
mc-healthcheck status-bedrock --host localhost --port 19132 --timeout 5s
mc-healthcheck status-bedrock --transport nethernet --tls --host play.example.com
mc-healthcheck status-bedrock --transport raknet
`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		host, _ := cmd.Flags().GetString("host")
		port, _ := cmd.Flags().GetInt("port")
		timeout, _ := cmd.Flags().GetDuration("timeout")
		transport, _ := cmd.Flags().GetString("transport")
		useTLS, _ := cmd.Flags().GetBool("tls")

		if port < 1 || port > 65535 {
			return fmt.Errorf("port must be between 1 and 65535, got %d", port)
		}

		if timeout <= 0 {
			return fmt.Errorf("timeout must be positive, got %s", timeout)
		}

		switch transport {
		case transportNetherNet:
			return checkNetherNet(host, port, timeout, useTLS)
		case transportRakNet:
			return checkRakNet(host, port, timeout)
		case transportAuto:
			nnErr := checkNetherNet(host, port, timeout, useTLS)
			if nnErr == nil {
				return nil
			}
			rnErr := checkRakNet(host, port, timeout)
			if rnErr != nil {
				return errors.Join(nnErr, rnErr)
			}
			_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "warning: nethernet check failed, raknet fallback succeeded: %v\n", nnErr)
			return nil
		default:
			return fmt.Errorf("transport must be %q, %q or %q, got %q", transportAuto, transportNetherNet, transportRakNet, transport)
		}
	},
}

func checkNetherNet(host string, port int, timeout time.Duration, useTLS bool) error {
	if err := nethernet.Status(host, port, timeout, useTLS); err != nil {
		return fmt.Errorf("status-bedrock nethernet check failed: %w", err)
	}
	return nil
}

func checkRakNet(host string, port int, timeout time.Duration) error {
	addr := net.JoinHostPort(host, strconv.Itoa(port))
	response, err := raknet.PingTimeout(addr, timeout)
	if err != nil {
		return fmt.Errorf("status-bedrock raknet check failed: %w", err)
	}
	if len(response) == 0 {
		return fmt.Errorf("status-bedrock raknet check failed: empty response from %s", addr)
	}
	return nil
}

func init() {
	statusBedrockCmd.Flags().String("host", "localhost", "server hostname")
	statusBedrockCmd.Flags().Int("port", 19132, "server port (TCP for NetherNet signaling, UDP for RakNet)")
	statusBedrockCmd.Flags().Duration("timeout", 5*time.Second, "connection and read timeout")
	statusBedrockCmd.Flags().String("transport", transportAuto, "auto (NetherNet, then RakNet fallback), nethernet or raknet")
	statusBedrockCmd.Flags().Bool("tls", false, "use https for the NetherNet signaling request")
	RootCmd.AddCommand(statusBedrockCmd)
}
