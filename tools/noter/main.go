// noter — standalone note manager
// Built into bin/noter by `termbox build`.
// Uses $NOTE_FOLDER (set in termbox.env) for storage.
//
// Usage:
//   noter <name>      open or create <name>.md in $EDITOR
//   noter list        list all notes with modification timestamps
//   noter search      fuzzy search notes (requires fzf + bat)
//   noter dir         print the active notes directory
package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "noter: %v\n", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	notesDir := resolveNotesDir()
	if err := os.MkdirAll(notesDir, 0755); err != nil {
		return fmt.Errorf("creating notes dir %q: %w", notesDir, err)
	}

	if len(args) == 0 {
		return listNotes(notesDir)
	}

	switch args[0] {
	case "list", "ls":
		return listNotes(notesDir)
	case "search", "s":
		return searchNotes(notesDir)
	case "dir":
		fmt.Println(notesDir)
		return nil
	default:
		return openNote(notesDir, args[0])
	}
}

func resolveNotesDir() string {
	if d := os.Getenv("NOTE_FOLDER"); d != "" {
		return d
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, "Developer", "notes")
}

func openNote(dir, name string) error {
	if !strings.HasSuffix(name, ".md") {
		name += ".md"
	}
	notePath := filepath.Join(dir, name)

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

func listNotes(dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("reading notes: %w", err)
	}

	var notes []os.DirEntry
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".md") {
			notes = append(notes, e)
		}
	}

	if len(notes) == 0 {
		fmt.Printf("  (no notes yet in %s)\n", dir)
		fmt.Println("  Create one with: noter <name>")
		return nil
	}

	fmt.Printf("\n  Notes — %s\n\n", dir)
	for _, e := range notes {
		info, _ := e.Info()
		name := strings.TrimSuffix(e.Name(), ".md")
		mod := info.ModTime().Format("2006-01-02 15:04")
		fmt.Printf("  %-36s  %s\n", name, mod)
	}
	fmt.Println()
	return nil
}

func searchNotes(dir string) error {
	if _, err := exec.LookPath("fzf"); err != nil {
		return fmt.Errorf("fzf is required for search")
	}
	preview := "bat --color=always {}"
	if _, err := exec.LookPath("bat"); err != nil {
		preview = "cat {}"
	}

	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "nvim"
	}

	c := exec.Command("fzf",
		"--preview", preview,
		"--height=80%", "--reverse", "--border",
		"--bind", fmt.Sprintf("enter:execute(%s {})", editor),
		"--header", fmt.Sprintf("Notes (%s) — Enter to open", time.Now().Format("2006-01-02")),
	)
	c.Dir = dir
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	c.Stdin = os.Stdin
	if err := c.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 130 {
			return nil // fzf cancelled
		}
		return err
	}
	return nil
}
