// Package cmd wires up the mc-healthcheck CLI: two subcommands, status and
// status-bedrock, each a thin wrapper around a protocol client.
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// RootCmd represents the base command when called without any subcommands.
var RootCmd = &cobra.Command{
	Use:   "mc-healthcheck",
	Short: "Minimal health checks for Minecraft Java and Bedrock servers",
	Long: `
mc-healthcheck checks whether a Minecraft server is responding, for use in
container health checks. It has exactly two subcommands:

  status          Java Edition Server List Ping (TCP)
  status-bedrock  Bedrock Edition RakNet unconnected ping (UDP)

Each exits 0 on a well-formed response within the timeout, non-zero otherwise.
`,
	SilenceUsage:  true,
	SilenceErrors: true,
	// Invoking mc-healthcheck with no subcommand must not exit 0: a
	// container HEALTHCHECK that forgets to specify status/status-bedrock
	// should fail loudly rather than pass by accident.
	RunE: func(cmd *cobra.Command, _ []string) error {
		_ = cmd.Help()
		return fmt.Errorf("no subcommand specified")
	},
}

// Execute adds all child commands to the root command and executes it. This
// is called by main.main(). It only needs to happen once to RootCmd.
func Execute() error {
	return RootCmd.Execute()
}
