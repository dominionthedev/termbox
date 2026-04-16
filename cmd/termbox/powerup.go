package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/dominionthedev/termbox/internal/envutil"
	"github.com/dominionthedev/termbox/internal/powerup"
	"github.com/dominionthedev/termbox/internal/registry"
	"github.com/spf13/cobra"
)

var powerupCmd = &cobra.Command{
	Use:   "powerup [list|activate|info|detect|run]",
	Short: "Manage powerups — purpose-specific packs of scripts and tools",
	Long: `Powerups are curated script+tool packs for a specific workflow.

  termbox powerup list               list available powerups with status
  termbox powerup activate <n>    activate a powerup (writes to termbox.env)
  termbox powerup info <n>        show contents and rules of a powerup
  termbox powerup detect             auto-activate powerups whose env criteria match
  termbox powerup run <pu> <script>  run a script from a specific powerup`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		action := args[0]
		home, err := registry.FindHome()
		if err != nil {
			return err
		}
		switch action {
		case "list":
			return powerupList(home)
		case "activate":
			if len(args) < 2 {
				return fmt.Errorf("usage: termbox powerup activate <n>")
			}
			return powerupActivate(home, args[1])
		case "info":
			if len(args) < 2 {
				return fmt.Errorf("usage: termbox powerup info <n>")
			}
			return powerupInfo(home, args[1])
		case "detect":
			return powerupDetect(home)
		case "run":
			if len(args) < 3 {
				return fmt.Errorf("usage: termbox powerup run <powerup> <script> [args...]")
			}
			return powerupRun(home, args[1], args[2], args[3:])
		default:
			return fmt.Errorf("unknown action %q — use: list, activate, info, detect, run", action)
		}
	},
}

func powerupList(home string) error {
	all, err := powerup.LoadAll(home)
	if err != nil {
		fmt.Println("  (no powerups directory)")
		return nil
	}
	active := os.Getenv("TERMBOX_POWERUP")
	fmt.Println("\n  Powerups\n")
	for _, p := range all {
		marker := "  "
		if p.Name == active {
			marker = "▶ "
		}

		ok, missing := p.MeetsRequires()
		meetsEnv := p.MeetsEnvCriteria()

		status := ""
		if !ok {
			status = fmt.Sprintf("  (needs: %s)", missing)
		} else if !meetsEnv && len(p.Rules.Env) > 0 {
			status = "  (env criteria not met)"
		}

		fmt.Printf("  %s%-16s  v%-8s  %s%s\n", marker, p.Name, p.Version, p.Description, status)
	}
	fmt.Println()
	return nil
}

func powerupActivate(home, name string) error {
	p, err := powerup.Load(home, name)
	if err != nil {
		return fmt.Errorf("powerup %q not found: %w", name, err)
	}

	ok, missing := p.MeetsRequires()
	if !ok {
		fmt.Fprintf(os.Stderr, "  ⚠  powerup %q requires %q (not found in PATH)\n", name, missing)
		fmt.Fprintf(os.Stderr, "     Activating anyway — install %q to use its scripts\n", missing)
	}

	envPath := filepath.Join(home, "config", "termbox.env")
	if err := envutil.UpdateVar(envPath, "TERMBOX_POWERUP", name); err != nil {
		return err
	}
	fmt.Printf("  ✓ Powerup activated: %s\n  Reload shell: source ~/.zshrc\n", name)
	return nil
}

func powerupInfo(home, name string) error {
	p, err := powerup.Load(home, name)
	if err != nil {
		return fmt.Errorf("powerup %q not found: %w", name, err)
	}
	fmt.Printf("\n  %s  v%s\n  %s\n\n", p.Name, p.Version, p.Description)
	if len(p.Scripts) > 0 {
		fmt.Printf("  Scripts:   %s\n", strings.Join(p.Scripts, ", "))
	}
	if len(p.Tools) > 0 {
		fmt.Printf("  Tools:     %s\n", strings.Join(p.Tools, ", "))
	}
	if len(p.Rules.Requires) > 0 {
		fmt.Printf("  Requires:  %s\n", strings.Join(p.Rules.Requires, ", "))
	}
	if len(p.Rules.Env) > 0 {
		fmt.Printf("  Env rules:\n")
		for _, rule := range p.Rules.Env {
			fmt.Printf("    %-14s  %s\n", rule.Kind, rule.Value)
		}
	}
	fmt.Println()
	return nil
}

func powerupDetect(home string) error {
	all, err := powerup.LoadAll(home)
	if err != nil {
		return fmt.Errorf("loading powerups: %w", err)
	}
	fmt.Println("\n  Auto-detect results (current directory):\n")
	for _, p := range all {
		if p.ShouldAutoActivate() {
			fmt.Printf("  ✓ %-16s  would activate\n", p.Name)
		} else {
			fmt.Printf("  ✗ %-16s  criteria not met\n", p.Name)
		}
	}
	fmt.Println()
	return nil
}

func powerupRun(home, puName, scriptName string, args []string) error {
	p, err := powerup.Load(home, puName)
	if err != nil {
		return fmt.Errorf("powerup %q not found: %w", puName, err)
	}

	found := false
	for _, s := range p.Scripts {
		if s == scriptName {
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("script %q is not in powerup %q", scriptName, puName)
	}

	// Delegate to the run command
	runArgs := append([]string{scriptName}, args...)
	return runCmd.RunE(runCmd, runArgs)
}

func init() {
	rootCmd.AddCommand(powerupCmd)
}
