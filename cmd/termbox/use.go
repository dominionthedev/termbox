package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/dominionthedev/termbox/internal/registry"
	"github.com/spf13/cobra"
)

var useCmd = &cobra.Command{
	Use:   "use [app[.name]] [destination]",
	Short: "Apply a config to its target app (copies, never symlinks)",
	Long: `Apply a termbox config to its target application.

  termbox use nvim                 list all config variants for nvim
  termbox use nvim.default         copy config/nvim/default.nvim/ → ~/.config/nvim/
  termbox use alacritty.cyberpunk  copy config/alacritty/cyberpunk.toml → registered target
  termbox use alacritty.cyberpunk ~/alt/alacritty.toml  copy to a custom destination

The current config at the destination is always backed up into
config/<app>/backup_<timestamp>/ before being replaced.

Addition configs (like zsh) print sourcing instructions instead.`,
	Args: cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		arg := args[0]

		// Optional destination override
		destOverride := ""
		if len(args) == 2 {
			destOverride = registry.ExpandHome(args[1])
		}

		home, err := registry.FindHome()
		if err != nil {
			return err
		}
		reg, err := registry.LoadRegistry(cfgFile)
		if err != nil {
			return fmt.Errorf("loading registry: %w", err)
		}

		if strings.Contains(arg, ".") {
			parts := strings.SplitN(arg, ".", 2)
			return applyConfig(home, reg, parts[0], parts[1], destOverride)
		}
		return listAppConfigs(reg, arg)
	},
}

func applyConfig(home string, reg *registry.Registry, app, name, destOverride string) error {
	var found *registry.Item
	for i := range reg.Configs {
		c := &reg.Configs[i]
		if strings.EqualFold(c.App, app) && c.Name == name {
			found = c
			break
		}
	}
	if found == nil {
		return fmt.Errorf("config %q not found for app %q — run 'termbox use %s' to list", name, app, app)
	}

	if found.ConfigKind == registry.KindAddition {
		fmt.Printf("  ℹ  %s.%s is an addition config — it's sourced from your shell.\n", app, name)
		fmt.Printf("  Run 'termbox setup' for the lines to add to ~/.zshrc.\n")
		return nil
	}

	// Resolve destination: override > registered target > error
	var dst string
	switch {
	case destOverride != "":
		dst = destOverride
	case found.Target != "":
		dst = registry.ExpandHome(found.Target)
	default:
		return fmt.Errorf("config %q has no target — add a target field in registry.yaml, or pass a destination path", found.Name)
	}

	srcPath := filepath.Join(home, found.Path)
	if _, err := os.Stat(srcPath); err != nil {
		return fmt.Errorf("source %q not found: %w", srcPath, err)
	}

	// Backup what's currently at dst
	if _, err := os.Stat(dst); err == nil {
		ts := time.Now().Format("20060102_150405")
		backupDst := filepath.Join(home, "config", app, fmt.Sprintf("backup_%s", ts))
		fmt.Printf("  ⤷ Backing up %s → config/%s/backup_%s\n", dst, app, ts)
		if err := copyPath(dst, backupDst); err != nil {
			return fmt.Errorf("backing up: %w", err)
		}
		if err := os.RemoveAll(dst); err != nil {
			return fmt.Errorf("clearing destination: %w", err)
		}
	}

	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return fmt.Errorf("creating parent dir: %w", err)
	}
	if err := copyPath(srcPath, dst); err != nil {
		return fmt.Errorf("copying config: %w", err)
	}

	fmt.Printf("  ✓ Applied %s.%s → %s\n", app, name, dst)
	return nil
}

func listAppConfigs(reg *registry.Registry, app string) error {
	configs := reg.ConfigsForApp(app)
	if len(configs) == 0 {
		return fmt.Errorf("no configs registered for app %q", app)
	}
	fmt.Printf("\n  Configs for %s:\n\n", app)
	for _, c := range configs {
		marker := "  "
		if c.Active {
			marker = "▶ "
		}
		target := c.Target
		if target == "" {
			target = "(no target — " + string(c.ConfigKind) + ")"
		}
		fmt.Printf("  %s%-20s  %-12s  %s\n", marker, c.Name, string(c.ConfigKind), c.Description)
		if c.ConfigKind != registry.KindAddition && c.Target != "" {
			fmt.Printf("      → %s\n", c.Target)
		}
	}
	fmt.Printf("\n  Pass a destination to override: termbox use %s.<n> <path>\n\n", app)
	return nil
}

// ── file copy helpers ─────────────────────────────────────────────────────────

func copyPath(src, dst string) error {
	info, err := os.Stat(src)
	if err != nil {
		return err
	}
	if info.IsDir() {
		return copyDir(src, dst)
	}
	return copyFile(src, dst)
}

func copyDir(src, dst string) error {
	if err := os.MkdirAll(dst, 0755); err != nil {
		return err
	}
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}
	for _, e := range entries {
		if err := copyPath(filepath.Join(src, e.Name()), filepath.Join(dst, e.Name())); err != nil {
			return err
		}
	}
	return nil
}

func copyFile(src, dst string) error {
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	if _, err := io.Copy(out, in); err != nil {
		return err
	}
	info, _ := os.Stat(src)
	if info != nil {
		_ = os.Chmod(dst, info.Mode())
	}
	return nil
}

func init() {
	rootCmd.AddCommand(useCmd)
}
