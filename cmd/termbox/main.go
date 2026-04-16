package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "termbox",
	Short: "Termbox — your terminal environment, organised",
	Long: `Termbox manages your terminal tools, scripts, configs, and powerups.

First time? Run:
  termbox setup --wizard

Then add exactly these two lines to the top of your ~/.zshrc:

  # ---- termbox env ----
  source $HOME/Developer/termbox/config/termbox.env
  # ---- termbox zsh ----
  source $TERMBOX_HOME/config/shell/zshrc

Termbox is a guest in your shell — it never edits your dotfiles.`,
	SilenceErrors: true,
	SilenceUsage:  true,
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "",
		"registry file (default: $TERMBOX_HOME/config/registry.yaml)")
}
