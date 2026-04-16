package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/dominionthedev/termbox/internal/registry"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run <script> [args...]",
	Short: "Execute a registered script by name",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		scriptArgs := args[1:]

		reg, err := registry.LoadRegistry(cfgFile)
		if err != nil {
			return fmt.Errorf("loading registry: %w", err)
		}
		item := reg.FindItem(name)
		if item == nil {
			return fmt.Errorf("script %q not found — run 'termbox list' to see available", name)
		}
		if item.Kind != "script" {
			return fmt.Errorf("%q is a %q, not a script", name, item.Kind)
		}

		home, err := registry.FindHome()
		if err != nil {
			return err
		}
		scriptPath := filepath.Join(home, item.Path)
		if _, err := os.Stat(scriptPath); err != nil {
			return fmt.Errorf("script not found at %q: %w", scriptPath, err)
		}
		if err := os.Chmod(scriptPath, 0755); err != nil {
			fmt.Fprintf(os.Stderr, "  ⚠  chmod: %v\n", err)
		}

		shell := os.Getenv("SHELL")
		if shell == "" {
			shell = "sh"
		}
		c := exec.Command(shell, append([]string{scriptPath}, scriptArgs...)...)
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		c.Stdin = os.Stdin
		c.Dir = home
		if err := c.Run(); err != nil {
			if e, ok := err.(*exec.ExitError); ok {
				os.Exit(e.ExitCode())
			}
			return fmt.Errorf("running %q: %w", name, err)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
