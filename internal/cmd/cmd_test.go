package cmd

import (
	"bytes"
	"strings"
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
	err := execute("status", "mc.example.com")
	if err == nil {
		t.Fatal("Execute() = nil, want error for unexpected positional argument")
	}
	if want := `unknown command "mc.example.com"`; !strings.Contains(err.Error(), want) {
		t.Fatalf("Execute() error = %q, want it to contain %q", err, want)
	}
}

func TestStatusBedrock_RejectsPositionalArgs(t *testing.T) {
	err := execute("status-bedrock", "mc.example.com")
	if err == nil {
		t.Fatal("Execute() = nil, want error for unexpected positional argument")
	}
	if want := `unknown command "mc.example.com"`; !strings.Contains(err.Error(), want) {
		t.Fatalf("Execute() error = %q, want it to contain %q", err, want)
	}
}

func TestStatus_RejectsInvalidPort(t *testing.T) {
	err := execute("status", "--port", "70000")
	if err == nil {
		t.Fatal("Execute() = nil, want error for out-of-range port")
	}
	if want := "port must be between 1 and 65535"; !strings.Contains(err.Error(), want) {
		t.Fatalf("Execute() error = %q, want it to contain %q", err, want)
	}
}

func TestStatusBedrock_RejectsInvalidPort(t *testing.T) {
	err := execute("status-bedrock", "--port", "-1")
	if err == nil {
		t.Fatal("Execute() = nil, want error for out-of-range port")
	}
	if want := "port must be between 1 and 65535"; !strings.Contains(err.Error(), want) {
		t.Fatalf("Execute() error = %q, want it to contain %q", err, want)
	}
}
