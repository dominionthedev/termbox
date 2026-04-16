package main

import (
	"fmt"

	"github.com/dominionthedev/termbox/internal/registry"
	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:   "remove <name>",
	Short: "Unregister a component from the registry",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		reg, err := registry.LoadRegistry(cfgFile)
		if err != nil {
			return fmt.Errorf("loading registry: %w", err)
		}
		if reg.FindItem(name) == nil {
			return fmt.Errorf("item %q not found in registry", name)
		}
		filter := func(items []registry.Item) []registry.Item {
			var out []registry.Item
			for _, i := range items {
				if i.Name != name { out = append(out, i) }
			}
			return out
		}
		reg.Tools = filter(reg.Tools)
		reg.Scripts = filter(reg.Scripts)
		reg.Configs = filter(reg.Configs)
		if err := registry.SaveRegistry(reg, cfgFile); err != nil {
			return fmt.Errorf("saving: %w", err)
		}
		fmt.Printf("  ✓ Removed %q from registry\n", name)
		return nil
	},
}

func init() { rootCmd.AddCommand(removeCmd) }
