// spicer — spice loader for shell startup
// Built into bin/spicer by `termbox build`.
// A "spice" is a .sh file containing shell functions/aliases, with a
// metadata header. Spicer manages which spices are enabled and generates
// a loader.sh that your shell sources at startup.
//
// Spice file format:
//   # spicer:name      git-extras
//   # spicer:desc      Handy git aliases and functions
//   # spicer:tags      git,productivity
//   # spicer:version   1.0.0
//   alias glog='git log --oneline --graph'
//   gclean() { ... }
//
// Shell integration (add to zshrc after termbox zshrc):
//   [[ -f ~/.config/spicer/loader.sh ]] && source ~/.config/spicer/loader.sh
//
// Usage:
//   spicer list                list all installed spices
//   spicer add <file.sh>       install a spice
//   spicer enable <name>       enable a spice
//   spicer disable <name>      disable a spice (stays installed)
//   spicer remove <name>       remove a spice
//   spicer reload              regenerate loader.sh from enabled spices
//   spicer new <name>          scaffold a new spice file
package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const (
	spiceDir    = ".config/spicer/spices"
	registryFile = ".config/spicer/registry.json"
	loaderFile  = ".config/spicer/loader.sh"
)

type Spice struct {
	Name    string `json:"name"`
	Desc    string `json:"desc"`
	Tags    string `json:"tags"`
	Version string `json:"version"`
	File    string `json:"file"`
	Enabled bool   `json:"enabled"`
}

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "spicer: %v\n", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	home, _ := os.UserHomeDir()

	if len(args) == 0 {
		return listSpices(home)
	}

	switch args[0] {
	case "list", "ls":
		return listSpices(home)
	case "add":
		if len(args) < 2 {
			return fmt.Errorf("usage: spicer add <file.sh>")
		}
		return addSpice(home, args[1])
	case "enable":
		if len(args) < 2 {
			return fmt.Errorf("usage: spicer enable <name>")
		}
		return toggleSpice(home, args[1], true)
	case "disable":
		if len(args) < 2 {
			return fmt.Errorf("usage: spicer disable <name>")
		}
		return toggleSpice(home, args[1], false)
	case "remove":
		if len(args) < 2 {
			return fmt.Errorf("usage: spicer remove <name>")
		}
		return removeSpice(home, args[1])
	case "reload":
		return generateLoader(home)
	case "new":
		if len(args) < 2 {
			return fmt.Errorf("usage: spicer new <name>")
		}
		return newSpice(home, args[1])
	default:
		return fmt.Errorf("unknown command %q — use: list, add, enable, disable, remove, reload, new", args[0])
	}
}

// ── Registry ─────────────────────────────────────────────────────────────────

func loadRegistry(home string) ([]Spice, error) {
	path := filepath.Join(home, registryFile)
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return []Spice{}, nil
	}
	if err != nil {
		return nil, err
	}
	var spices []Spice
	return spices, json.Unmarshal(data, &spices)
}

func saveRegistry(home string, spices []Spice) error {
	path := filepath.Join(home, registryFile)
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(spices, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// ── Commands ──────────────────────────────────────────────────────────────────

func listSpices(home string) error {
	spices, err := loadRegistry(home)
	if err != nil {
		return err
	}
	if len(spices) == 0 {
		fmt.Printf("  (no spices installed — add one with: spicer add <file.sh>)\n")
		return nil
	}
	fmt.Printf("\n  Spices — %s\n\n", filepath.Join(home, spiceDir))
	for _, s := range spices {
		marker := "  "
		if s.Enabled {
			marker = "▶ "
		}
		fmt.Printf("  %s%-20s  v%-8s  %s\n", marker, s.Name, s.Version, s.Desc)
	}
	fmt.Println()
	return nil
}

func addSpice(home, srcFile string) error {
	spice, err := parseSpiceHeader(srcFile)
	if err != nil {
		return err
	}
	if spice.Name == "" {
		return fmt.Errorf("spice file must have a '# spicer:name' header line")
	}

	spices, err := loadRegistry(home)
	if err != nil {
		return err
	}
	for _, s := range spices {
		if s.Name == spice.Name {
			return fmt.Errorf("spice %q already installed — remove it first", spice.Name)
		}
	}

	// Copy spice file to spice dir
	destDir := filepath.Join(home, spiceDir)
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return err
	}
	destFile := filepath.Join(destDir, spice.Name+".sh")
	if err := copyFile(srcFile, destFile); err != nil {
		return err
	}
	spice.File = destFile
	spice.Enabled = true

	spices = append(spices, *spice)
	if err := saveRegistry(home, spices); err != nil {
		return err
	}
	if err := generateLoader(home); err != nil {
		return err
	}

	fmt.Printf("  ✓ Added spice: %s\n", spice.Name)
	fmt.Printf("  Reload shell to activate: source ~/.zshrc\n")
	return nil
}

func toggleSpice(home, name string, enable bool) error {
	spices, err := loadRegistry(home)
	if err != nil {
		return err
	}
	found := false
	for i := range spices {
		if spices[i].Name == name {
			spices[i].Enabled = enable
			found = true
		}
	}
	if !found {
		return fmt.Errorf("spice %q not found", name)
	}
	if err := saveRegistry(home, spices); err != nil {
		return err
	}
	if err := generateLoader(home); err != nil {
		return err
	}
	verb := "enabled"
	if !enable {
		verb = "disabled"
	}
	fmt.Printf("  ✓ Spice %s: %s\n  Reload shell: source ~/.zshrc\n", verb, name)
	return nil
}

func removeSpice(home, name string) error {
	spices, err := loadRegistry(home)
	if err != nil {
		return err
	}
	var kept []Spice
	removed := false
	for _, s := range spices {
		if s.Name == name {
			_ = os.Remove(s.File)
			removed = true
		} else {
			kept = append(kept, s)
		}
	}
	if !removed {
		return fmt.Errorf("spice %q not found", name)
	}
	if err := saveRegistry(home, kept); err != nil {
		return err
	}
	_ = generateLoader(home)
	fmt.Printf("  ✓ Removed spice: %s\n", name)
	return nil
}

func generateLoader(home string) error {
	spices, err := loadRegistry(home)
	if err != nil {
		return err
	}
	path := filepath.Join(home, loaderFile)
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	var sb strings.Builder
	sb.WriteString("# spicer loader.sh — generated by spicer reload\n")
	sb.WriteString("# Source this from your zshrc: [[ -f ~/.config/spicer/loader.sh ]] && source ~/.config/spicer/loader.sh\n\n")

	count := 0
	for _, s := range spices {
		if !s.Enabled {
			continue
		}
		if _, err := os.Stat(s.File); err != nil {
			continue
		}
		sb.WriteString(fmt.Sprintf("# spice: %s\n", s.Name))
		sb.WriteString(fmt.Sprintf("source %q\n\n", s.File))
		count++
	}

	if err := os.WriteFile(path, []byte(sb.String()), 0644); err != nil {
		return err
	}
	fmt.Printf("  ✓ loader.sh regenerated (%d spice(s) active)\n", count)
	return nil
}

func newSpice(home, name string) error {
	destDir := filepath.Join(home, spiceDir)
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return err
	}
	path := filepath.Join(destDir, name+".sh")
	if _, err := os.Stat(path); err == nil {
		return fmt.Errorf("spice file already exists at %s", path)
	}

	scaffold := fmt.Sprintf(`#!/usr/bin/env bash
# spicer:name      %s
# spicer:desc      Description of what this spice does
# spicer:tags      tag1,tag2
# spicer:version   0.1.0

# Add your shell functions and aliases below.
# They will be loaded into your shell at startup via spicer.

# Example:
# alias example='echo "hello from %s"'
# example_func() { echo "$1"; }
`, name, name)

	if err := os.WriteFile(path, []byte(scaffold), 0755); err != nil {
		return err
	}
	fmt.Printf("  ✓ Created: %s\n", path)
	fmt.Printf("  Edit it, then install with: spicer add %s\n", path)
	return nil
}

// ── Helpers ───────────────────────────────────────────────────────────────────

func parseSpiceHeader(file string) (*Spice, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, fmt.Errorf("opening spice file: %w", err)
	}
	defer f.Close()

	spice := &Spice{}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if !strings.HasPrefix(line, "# spicer:") {
			if spice.Name != "" {
				break // stop after the header block ends
			}
			continue
		}
		rest := strings.TrimPrefix(line, "# spicer:")
		parts := strings.SplitN(rest, " ", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])
		switch key {
		case "name":    spice.Name = val
		case "desc":    spice.Desc = val
		case "tags":    spice.Tags = val
		case "version": spice.Version = val
		}
	}
	return spice, scanner.Err()
}

func copyFile(src, dst string) error {
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
	_, err = io.Copy(out, in)
	return err
}
