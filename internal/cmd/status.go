package cmd

import (
	"encoding/json"
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
mc-healthcheck status --host localhost --port 25565 --json
`,
	RunE: func(cmd *cobra.Command, _ []string) error {
		host, _ := cmd.Flags().GetString("host")
		port, _ := cmd.Flags().GetInt("port")
		timeout, _ := cmd.Flags().GetDuration("timeout")
		printJSON, _ := cmd.Flags().GetBool("json")

		if timeout <= 0 {
			return fmt.Errorf("timeout must be positive, got %s", timeout)
		}

		if !printJSON {
			if err := slp.Status(host, port, timeout); err != nil {
				return fmt.Errorf("status check failed: %w", err)
			}
			return nil
		}

		payload, err := slp.Query(host, port, timeout)
		if err != nil {
			return fmt.Errorf("status check failed: %w", err)
		}
		encoded, err := json.Marshal(payload)
		if err != nil {
			return fmt.Errorf("encode status response: %w", err)
		}
		if _, err := fmt.Fprintln(cmd.OutOrStdout(), string(encoded)); err != nil {
			return fmt.Errorf("write status response: %w", err)
		}
		return nil
	},
}

func init() {
	statusCmd.Flags().String("host", "localhost", "server hostname")
	statusCmd.Flags().Int("port", 25565, "server port")
	statusCmd.Flags().Duration("timeout", 5*time.Second, "connection and read timeout")
	statusCmd.Flags().Bool("json", false, "print the decoded status response as JSON to stdout on success")
	RootCmd.AddCommand(statusCmd)
}
