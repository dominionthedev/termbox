package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/dominionthedev/termbox/internal/registry"
	"github.com/dominionthedev/termbox/internal/settings"
	"github.com/spf13/cobra"
)

var buildCmd = &cobra.Command{
	Use:   "build [name...]",
	Short: "Compile registered tools into bin/",
	Long: `Compile all registered tools, or specific ones by name.

  termbox build              build all tools
  termbox build termbox      build only the termbox CLI
  termbox build noter banner build specific tools`,
	RunE: func(cmd *cobra.Command, args []string) error {
		reg, err := registry.LoadRegistry(cfgFile)
		if err != nil {
			return fmt.Errorf("loading registry: %w", err)
		}
		home, err := registry.FindHome()
		if err != nil {
			return err
		}

		cfg, _ := settings.Load(home)
		outDir := cfg.Build.OutputDir
		if outDir == "" {
			outDir = "bin"
		}
		binDir := filepath.Join(home, outDir)
		if err := os.MkdirAll(binDir, 0755); err != nil {
			return fmt.Errorf("creating %s: %w", binDir, err)
		}

		// Filter to requested tools only (or all if no args)
		filter := map[string]bool{}
		for _, a := range args {
			filter[a] = true
		}

		fmt.Printf("\n  Building tools → %s/\n\n", outDir)
		built, failed := 0, 0

		for _, t := range reg.Tools {
			if len(filter) > 0 && !filter[t.Name] {
				continue
			}
			if t.Kind != "tool" || t.Path == "" {
				continue
			}

			out := filepath.Join(binDir, t.Name)
			srcPath := "./" + t.Path

			extraFlags := []string{}
			if cfg.Build.GoFlags != "" {
				extraFlags = append(extraFlags, cfg.Build.GoFlags)
			}

			goBuildArgs := append([]string{"build", "-o", out}, extraFlags...)
			goBuildArgs = append(goBuildArgs, srcPath)

			c := exec.Command("go", goBuildArgs...)
			c.Dir = home
			c.Stdout = os.Stdout
			c.Stderr = os.Stderr

			fmt.Printf("  → %-16s  %s\n", t.Name, srcPath)
			if err := c.Run(); err != nil {
				fmt.Fprintf(os.Stderr, "  ✗ %s: %v\n", t.Name, err)
				failed++
			} else {
				fmt.Printf("  ✓ %-16s  %s/%s\n", t.Name, outDir, t.Name)
				built++
			}
		}

		if built == 0 && failed == 0 {
			fmt.Println("  (no tools to build)")
		}
		fmt.Printf("\n  %d built, %d failed\n\n", built, failed)
		if failed > 0 {
			return fmt.Errorf("%d build(s) failed", failed)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)
}
