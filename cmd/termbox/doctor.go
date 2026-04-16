package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/dominionthedev/termbox/internal/registry"
	"github.com/spf13/cobra"
)

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Health check your termbox environment",
	RunE: func(cmd *cobra.Command, args []string) error {
		pass, fail, warn := "  ✓", "  ✗", "  ⚠"
		fmt.Println("\n  ╭──────────────────────────────╮")
		fmt.Println("  │  Termbox Doctor              │")
		fmt.Println("  ╰──────────────────────────────╯")

		fmt.Println("\n  Termbox")
		home, err := registry.FindHome()
		if err != nil {
			fmt.Printf("%s TERMBOX_HOME not set — run: termbox setup --wizard\n", fail)
		} else {
			fmt.Printf("%s TERMBOX_HOME = %s\n", pass, home)
		}
		check := func(path, label string) {
			if _, err := os.Stat(path); err == nil {
				fmt.Printf("%s %s\n", pass, label)
			} else {
				fmt.Printf("%s %s (not found)\n", warn, label)
			}
		}
		check(filepath.Join(home, "config", "termbox.env"), "termbox.env")
		check(filepath.Join(home, "config", "theme.env"), "theme.env")
		check(filepath.Join(home, "config", "registry.yaml"), "registry.yaml")

		fmt.Println("\n  Core tools")
		for _, t := range []string{"tmux", "nvim", "starship", "fzf", "zoxide", "bat", "eza", "fd", "rg"} {
			if _, err := exec.LookPath(t); err == nil {
				fmt.Printf("%s %s\n", pass, t)
			} else {
				fmt.Printf("%s %s\n", fail, t)
			}
		}

		fmt.Println("\n  Optional tools")
		for _, t := range []string{"lolcat", "tokei", "podman", "docker", "lazygit", "delta", "procs", "dust", "htop"} {
			if _, err := exec.LookPath(t); err == nil {
				fmt.Printf("%s %s\n", pass, t)
			} else {
				fmt.Printf("%s %s\n", warn, t)
			}
		}
		fmt.Println()
		return nil
	},
}

func init() { rootCmd.AddCommand(doctorCmd) }
