package main

import (
	"fmt"
	"github.com/spf13/cobra"
)

// theme is a thin alias pointing to sheme — the real theme engine.
var themeCmd = &cobra.Command{
	Use:   "theme",
	Short: "Alias for 'termbox sheme' — use sheme for full theme control",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("  'termbox theme' → use 'termbox sheme' for full control:")
		fmt.Println()
		fmt.Println("    termbox sheme list")
		fmt.Println("    termbox sheme palettes")
		fmt.Println("    termbox sheme from <palette>")
		fmt.Println("    termbox sheme apply <n>")
		fmt.Println("    termbox sheme new <n>")
		fmt.Println("    termbox sheme info <n>")
		return nil
	},
}

func init() { rootCmd.AddCommand(themeCmd) }
