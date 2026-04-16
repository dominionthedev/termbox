// shelldoctor — Shell environment health checker
// Built into bin/shelldoctor by `termbox build`.
// Part of the termbox shell tools suite.
//
// Usage:
//   shelldoctor             run all checks
//   shelldoctor path        check PATH only
//   shelldoctor config      check shell config files
//   shelldoctor tools       check required tools
//   shelldoctor startup     benchmark shell startup time
package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const (
	pass = "  ✓"
	fail = "  ✗"
	warn = "  ⚠"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "shelldoctor: %v\n", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	check := "all"
	if len(args) > 0 {
		check = args[0]
	}

	fmt.Println()
	fmt.Println("  ╭──────────────────────────────────────╮")
	fmt.Println("  │  ShellDoctor                         │")
	fmt.Println("  ╰──────────────────────────────────────╯")

	switch check {
	case "all":
		checkPath()
		checkConfig()
		checkTools()
		checkStartup()
	case "path":
		checkPath()
	case "config":
		checkConfig()
	case "tools":
		checkTools()
	case "startup":
		checkStartup()
	default:
		return fmt.Errorf("unknown check %q — use: all, path, config, tools, startup", check)
	}
	return nil
}

func checkPath() {
	fmt.Println("\n  PATH")
	seen := map[string]int{}
	broken := 0

	pathEnv := os.Getenv("PATH")
	entries := filepath.SplitList(pathEnv)

	for _, entry := range entries {
		seen[entry]++
	}

	for _, entry := range entries {
		info, err := os.Stat(entry)
		switch {
		case err != nil:
			fmt.Printf("%s  %s  (does not exist)\n", fail, entry)
			broken++
		case !info.IsDir():
			fmt.Printf("%s  %s  (not a directory)\n", warn, entry)
		case seen[entry] > 1:
			fmt.Printf("%s  %s  (duplicate ×%d)\n", warn, entry, seen[entry])
		default:
			fmt.Printf("%s  %s\n", pass, entry)
		}
	}

	fmt.Printf("\n  %d entries, %d broken\n", len(entries), broken)
}

func checkConfig() {
	fmt.Println("\n  Shell Config")
	home, _ := os.UserHomeDir()

	configs := []struct{ name, path string }{
		{".zshrc", filepath.Join(home, ".zshrc")},
		{".zshenv", filepath.Join(home, ".zshenv")},
		{".zprofile", filepath.Join(home, ".zprofile")},
		{".bashrc", filepath.Join(home, ".bashrc")},
		{".bash_profile", filepath.Join(home, ".bash_profile")},
	}

	for _, cfg := range configs {
		info, err := os.Stat(cfg.path)
		if err != nil {
			continue // doesn't exist — not a problem
		}
		size := info.Size()
		if size == 0 {
			fmt.Printf("%s  %-20s  (empty)\n", warn, cfg.name)
			continue
		}
		data, _ := os.ReadFile(cfg.path)
		content := string(data)
		// Check for termbox integration
		hasTB := strings.Contains(content, "termbox.env") || strings.Contains(content, "TERMBOX_HOME")
		marker := pass
		note := ""
		if !hasTB && cfg.name == ".zshrc" {
			marker = warn
			note = "  (termbox not integrated — run setup.sh)"
		}
		fmt.Printf("%s  %-20s  %d bytes%s\n", marker, cfg.name, size, note)
	}
}

func checkTools() {
	fmt.Println("\n  Core tools")
	core := []string{"zsh", "nvim", "tmux", "git", "fzf", "bat", "eza", "fd", "rg", "starship", "zoxide"}
	for _, t := range core {
		if _, err := exec.LookPath(t); err == nil {
			fmt.Printf("%s  %s\n", pass, t)
		} else {
			fmt.Printf("%s  %s\n", fail, t)
		}
	}

	fmt.Println("\n  Optional tools")
	optional := []string{"lolcat", "delta", "lazygit", "podman", "docker", "tokei", "procs", "dust", "htop"}
	for _, t := range optional {
		if _, err := exec.LookPath(t); err == nil {
			fmt.Printf("%s  %s\n", pass, t)
		} else {
			fmt.Printf("%s  %s  (optional)\n", warn, t)
		}
	}
}

func checkStartup() {
	fmt.Println("\n  Startup time")
	shell := os.Getenv("SHELL")
	if shell == "" {
		fmt.Printf("%s  $SHELL not set\n", warn)
		return
	}

	start := time.Now()
	c := exec.Command(shell, "-i", "-c", "exit")
	c.Env = os.Environ()
	if err := c.Run(); err != nil {
		fmt.Printf("%s  Could not benchmark: %v\n", warn, err)
		return
	}
	elapsed := time.Since(start)

	status := pass
	note := ""
	switch {
	case elapsed > 1*time.Second:
		status = fail
		note = " (very slow — audit plugins)"
	case elapsed > 500*time.Millisecond:
		status = warn
		note = " (slow — consider lazy loading)"
	}

	fmt.Printf("%s  %s startup: %s%s\n", status, filepath.Base(shell), elapsed.Round(time.Millisecond), note)
}
