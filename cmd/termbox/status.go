package main

import (
	"fmt"
	"os"

	"github.com/dominionthedev/termbox/internal/registry"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show current termbox runtime status",
	RunE: func(cmd *cobra.Command, args []string) error {
		home, _ := registry.FindHome()
		fmt.Println()
		fmt.Printf("  TERMBOX_HOME         %s\n", home)
		fmt.Printf("  SHELL_THEME          %s\n", os.Getenv("SHELL_THEME"))
		fmt.Printf("  TERMBOX_POWERUP      %s\n", os.Getenv("TERMBOX_POWERUP"))
		fmt.Printf("  NOTE_FOLDER          %s\n", os.Getenv("NOTE_FOLDER"))
		fmt.Printf("  TERMBOX_DEFAULT_NVIM %s\n", os.Getenv("TERMBOX_DEFAULT_NVIM"))
		fmt.Printf("  TERMBOX_SHOW_BANNER  %s\n", os.Getenv("TERMBOX_SHOW_BANNER"))
		fmt.Println()
		return nil
	},
}

func init() { rootCmd.AddCommand(statusCmd) }
