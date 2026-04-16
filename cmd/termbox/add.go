package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/dominionthedev/termbox/internal/registry"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add <kind> <name> <path> <description>",
	Short: "Register a component in the registry",
	Long: `Kinds: tool | script | config
Use --app and --config-kind for configs.
Use --target for configs that are copied to a destination.`,
	Args: cobra.ExactArgs(4),
	RunE: func(cmd *cobra.Command, args []string) error {
		kind, name, path, desc := args[0], args[1], args[2], args[3]
		valid := map[string]bool{"tool": true, "script": true, "config": true}
		if !valid[kind] {
			return fmt.Errorf("unknown kind %q — use: tool, script, config", kind)
		}
		app, _ := cmd.Flags().GetString("app")
		ck, _ := cmd.Flags().GetString("config-kind")
		target, _ := cmd.Flags().GetString("target")

		reg, err := registry.LoadRegistry(cfgFile)
		if err != nil {
			return fmt.Errorf("loading registry: %w", err)
		}
		item := registry.Item{
			Name: name, Kind: kind, App: app,
			ConfigKind: registry.ConfigKind(ck),
			Path: path, Target: target, Description: desc,
		}
		if err := reg.ValidateNewItem(item); err != nil {
			if errors.Is(err, registry.ErrDuplicateName) {
				return fmt.Errorf("%w — use 'termbox remove %s' first", err, name)
			}
			return err
		}
		home, _ := registry.FindHome()
		if home != "" {
			if _, err := os.Stat(filepath.Join(home, path)); err != nil {
				fmt.Fprintf(os.Stderr, "  ⚠  path %q not found — registering anyway\n", path)
			}
		}
		switch kind {
		case "tool":    reg.Tools = append(reg.Tools, item)
		case "script":  reg.Scripts = append(reg.Scripts, item)
		case "config":  reg.Configs = append(reg.Configs, item)
		}
		if err := registry.SaveRegistry(reg, cfgFile); err != nil {
			return fmt.Errorf("saving: %w", err)
		}
		fmt.Printf("  ✓ Registered %s [%s]\n", name, kind)
		return nil
	},
}

func init() {
	addCmd.Flags().String("app", "", "app this config belongs to (nvim, starship...)")
	addCmd.Flags().String("config-kind", "main", "main | template | addition")
	addCmd.Flags().String("target", "", "copy destination for main/template configs")
	rootCmd.AddCommand(addCmd)
}
