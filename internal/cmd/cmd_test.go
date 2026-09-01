package cmd

import (
	"bytes"
	"testing"
)

func execute(args ...string) error {
	RootCmd.SetArgs(args)
	RootCmd.SetOut(&bytes.Buffer{})
	RootCmd.SetErr(&bytes.Buffer{})
	return RootCmd.Execute()
}

func TestRoot_NoSubcommandFails(t *testing.T) {
	if err := execute(); err == nil {
		t.Fatal("Execute() = nil, want error when no subcommand is given")
	}
}

func TestStatus_RejectsPositionalArgs(t *testing.T) {
	if err := execute("status", "mc.example.com"); err == nil {
		t.Fatal("Execute() = nil, want error for unexpected positional argument")
	}
}

func TestStatusBedrock_RejectsPositionalArgs(t *testing.T) {
	if err := execute("status-bedrock", "mc.example.com"); err == nil {
		t.Fatal("Execute() = nil, want error for unexpected positional argument")
	}
}

func TestStatus_RejectsInvalidPort(t *testing.T) {
	if err := execute("status", "--port", "70000"); err == nil {
		t.Fatal("Execute() = nil, want error for out-of-range port")
	}
}

func TestStatusBedrock_RejectsInvalidPort(t *testing.T) {
	if err := execute("status-bedrock", "--port", "-1"); err == nil {
		t.Fatal("Execute() = nil, want error for out-of-range port")
	}
}
