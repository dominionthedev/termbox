package main

import (
	"fmt"
	"os"
	"os/exec"
	"github.com/dominionthedev/termbox/internal/registry"
	"github.com/spf13/cobra"
)

var upgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "Pull latest changes and rebuild termbox",
	RunE: func(cmd *cobra.Command, args []string) error {
		home, err := registry.FindHome()
		if err != nil { return err }
		fmt.Println("  → Pulling latest changes...")
		pull := exec.Command("git", "-C", home, "pull")
		pull.Stdout = os.Stdout
		pull.Stderr = os.Stderr
		pull.Run()
		fmt.Println("  → Rebuilding termbox...")
		exe, _ := os.Executable()
		build := exec.Command("go", "build", "-o", exe, "./cmd/termbox")
		build.Dir = home
		build.Stdout = os.Stdout
		build.Stderr = os.Stderr
		if err := build.Run(); err != nil {
			return fmt.Errorf("build failed: %w", err)
		}
		fmt.Println("  ✓ Termbox upgraded.")
		return nil
	},
}

func init() { rootCmd.AddCommand(upgradeCmd) }
