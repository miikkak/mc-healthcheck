package cmd

import (
	"fmt"
	"time"

	"github.com/miikkak/mc-healthcheck/internal/slp"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check a Java Edition server via Server List Ping (TCP)",
	Example: `
mc-healthcheck status --host localhost --port 25565 --timeout 5s
`,
	RunE: func(cmd *cobra.Command, _ []string) error {
		host, _ := cmd.Flags().GetString("host")
		port, _ := cmd.Flags().GetInt("port")
		timeout, _ := cmd.Flags().GetDuration("timeout")

		if err := slp.Status(host, port, timeout); err != nil {
			return fmt.Errorf("status check failed: %w", err)
		}
		return nil
	},
}

func init() {
	statusCmd.Flags().String("host", "localhost", "server hostname")
	statusCmd.Flags().Int("port", 25565, "server port")
	statusCmd.Flags().Duration("timeout", 5*time.Second, "connection and read timeout")
	RootCmd.AddCommand(statusCmd)
}
