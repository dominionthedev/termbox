package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var noteCmd = &cobra.Command{
	Use:   "note [name|list|search]",
	Short: "Create or open notes (uses $NOTE_FOLDER from termbox.env)",
	Long: `Manage markdown notes in your configured notes folder.

  termbox note <n>      open or create <n>.md in $EDITOR
  termbox note list        list all notes with timestamps
  termbox note search      fuzzy search notes (requires fzf + bat)`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		action := args[0]

		notesDir := os.Getenv("NOTE_FOLDER")
		if notesDir == "" {
			home, _ := os.UserHomeDir()
			notesDir = filepath.Join(home, "Developer", "notes")
		}

		if err := os.MkdirAll(notesDir, 0755); err != nil {
			return fmt.Errorf("creating notes dir %q: %w", notesDir, err)
		}

		switch action {
		case "list":
			entries, err := os.ReadDir(notesDir)
			if err != nil {
				return fmt.Errorf("reading notes: %w", err)
			}
			if len(entries) == 0 {
				fmt.Printf("  (no notes yet in %s — create one with: termbox note <n>)\n", notesDir)
				return nil
			}
			fmt.Printf("\n  Notes in %s\n\n", notesDir)
			for _, e := range entries {
				if !e.IsDir() && strings.HasSuffix(e.Name(), ".md") {
					info, _ := e.Info()
					name := strings.TrimSuffix(e.Name(), ".md")
					fmt.Printf("  %-32s  %s\n", name, info.ModTime().Format("2006-01-02 15:04"))
				}
			}
			fmt.Println()

		case "search":
			if _, err := exec.LookPath("fzf"); err != nil {
				return fmt.Errorf("fzf is required for search — install it first")
			}
			preview := "bat --color=always {}"
			if _, err := exec.LookPath("bat"); err != nil {
				preview = "cat {}"
			}
			c := exec.Command("fzf",
				"--preview", preview,
				"--height=80%", "--reverse", "--border",
				"--bind", "enter:execute($EDITOR {})",
				"--header", "Notes — press Enter to open in $EDITOR",
			)
			c.Dir = notesDir
			c.Stdout = os.Stdout
			c.Stderr = os.Stderr
			c.Stdin = os.Stdin
			if err := c.Run(); err != nil {
				if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 130 {
					return nil // user cancelled fzf with Esc/Ctrl-C
				}
				return err
			}

		default:
			name := action
			if !strings.HasSuffix(name, ".md") {
				name += ".md"
			}
			notePath := filepath.Join(notesDir, name)

			editor := os.Getenv("EDITOR")
			if editor == "" {
				editor = "nvim"
			}
			c := exec.Command(editor, notePath)
			c.Stdout = os.Stdout
			c.Stderr = os.Stderr
			c.Stdin = os.Stdin
			return c.Run()
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(noteCmd)
}
