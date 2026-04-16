package main

import (
	"fmt"

	"github.com/dominionthedev/termbox/internal/registry"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all registered components",
	RunE: func(cmd *cobra.Command, args []string) error {
		reg, err := registry.LoadRegistry(cfgFile)
		if err != nil {
			return fmt.Errorf("loading registry: %w", err)
		}
		sec := func(title string, fn func()) {
			fmt.Printf("\n  %s\n  %s\n", title, "────────────────────────────────────────")
			fn()
		}
		sec("Tools", func() {
			for _, t := range reg.Tools {
				m := "  "; if t.Active { m = "▶ " }
				fmt.Printf("  %s%-20s  %s\n", m, t.Name, t.Description)
			}
		})
		sec("Scripts", func() {
			for _, s := range reg.Scripts {
				m := "  "; if s.Active { m = "▶ " }
				fmt.Printf("  %s%-20s  %s\n", m, s.Name, s.Description)
			}
		})
		sec("Configs", func() {
			for _, c := range reg.Configs {
				m := "  "; if c.Active { m = "▶ " }
				app := ""
				if c.App != "" { app = "[" + c.App + "] " }
				fmt.Printf("  %s%-20s  %s%s  (%s)\n", m, c.Name, app, c.Description, string(c.ConfigKind))
			}
		})
		fmt.Println()
		return nil
	},
}

func init() { rootCmd.AddCommand(listCmd) }
