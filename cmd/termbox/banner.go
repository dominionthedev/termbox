package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/dominionthedev/termbox/internal/registry"
	"github.com/dominionthedev/termbox/internal/envutil"
	"github.com/spf13/cobra"
)

var bannerCmd = &cobra.Command{
	Use:   "banner [show|set <file>|off|on]",
	Short: "Manage and display the startup banner",
	RunE: func(cmd *cobra.Command, args []string) error {
		action := "show"
		if len(args) > 0 {
			action = args[0]
		}

		home, err := registry.FindHome()
		if err != nil {
			return err
		}
		envPath := filepath.Join(home, "config", "termbox.env")

		switch action {
		case "show":
			return showBanner(home)
		case "set":
			if len(args) < 2 {
				return fmt.Errorf("usage: termbox banner set <file>")
			}
			if _, err := os.Stat(args[1]); err != nil {
				return fmt.Errorf("file not found: %q", args[1])
			}
			if err := envutil.UpdateVar(envPath, "BANNER", args[1]); err != nil {
				return err
			}
			fmt.Printf("  ✓ Banner file set to: %s\n  Reload shell to apply.\n", args[1])
			return nil
		case "off":
			if err := envutil.UpdateVar(envPath, "TERMBOX_SHOW_BANNER", "false"); err != nil {
				return err
			}
			fmt.Println("  ✓ Banner disabled. Reload shell: source ~/.zshrc")
			return nil
		case "on":
			if err := envutil.UpdateVar(envPath, "TERMBOX_SHOW_BANNER", "true"); err != nil {
				return err
			}
			fmt.Println("  ✓ Banner enabled. Reload shell: source ~/.zshrc")
			return nil
		default:
			return fmt.Errorf("unknown action %q — use: show, set <file>, on, off", action)
		}
	},
}

func showBanner(home string) error {
	bannerFile := os.Getenv("BANNER")
	if bannerFile == "" {
		bannerFile = filepath.Join(home, "assets", "dominiondev.banner")
	}

	data, err := os.ReadFile(bannerFile)
	if err != nil {
		return fmt.Errorf("banner file not found at %q\nSet it with: termbox banner set <file>", bannerFile)
	}

	if _, err := exec.LookPath("lolcat"); err == nil {
		r, w, _ := os.Pipe()
		go func() { w.Write(data); w.Close() }()
		c := exec.Command("lolcat")
		c.Stdin = r
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		if c.Run() == nil {
			return nil
		}
	}

	fmt.Print(string(data))
	return nil
}

func init() {
	rootCmd.AddCommand(bannerCmd)
}
