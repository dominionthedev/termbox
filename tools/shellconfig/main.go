// shellconfig — Shell configuration manager
// Built into bin/shellconfig by `termbox build`.
//
// Manages shell configuration files: view, backup, diff, and switch profiles.
// Does not modify configs directly — it works with copies and backups.
//
// Usage:
//   shellconfig view              view the current shell config
//   shellconfig backup            snapshot the current config with a timestamp
//   shellconfig restore [id]      restore a previous snapshot
//   shellconfig list backups      list available snapshots
//   shellconfig diff [a] [b]      diff two configs or snapshots
//   shellconfig detect            detect shell and print config file path
package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

const backupDir = ".config/shellconfig/backups"

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "shellconfig: %v\n", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	if len(args) == 0 {
		return detectShell()
	}

	switch args[0] {
	case "view":
		return viewConfig()
	case "backup":
		return backupConfig()
	case "restore":
		id := ""
		if len(args) > 1 {
			id = args[1]
		}
		return restoreConfig(id)
	case "list":
		return listBackups()
	case "diff":
		a, b := "", ""
		if len(args) > 1 { a = args[1] }
		if len(args) > 2 { b = args[2] }
		return diffConfigs(a, b)
	case "detect":
		return detectShell()
	default:
		return fmt.Errorf("unknown command %q — use: view, backup, restore, list, diff, detect", args[0])
	}
}

func detectShell() error {
	shell := os.Getenv("SHELL")
	if shell == "" {
		return fmt.Errorf("$SHELL is not set")
	}
	home, _ := os.UserHomeDir()

	var rcFile string
	switch filepath.Base(shell) {
	case "zsh":
		rcFile = filepath.Join(home, ".zshrc")
	case "bash":
		rcFile = filepath.Join(home, ".bashrc")
		if _, err := os.Stat(rcFile); err != nil {
			rcFile = filepath.Join(home, ".bash_profile")
		}
	case "fish":
		rcFile = filepath.Join(home, ".config", "fish", "config.fish")
	default:
		rcFile = filepath.Join(home, "."+filepath.Base(shell)+"rc")
	}

	fmt.Printf("  Shell:  %s\n", shell)
	fmt.Printf("  Config: %s\n", rcFile)

	if _, err := os.Stat(rcFile); err != nil {
		fmt.Printf("  Status: not found\n")
	} else {
		info, _ := os.Stat(rcFile)
		fmt.Printf("  Size:   %d bytes\n", info.Size())
		fmt.Printf("  Modified: %s\n", info.ModTime().Format("2006-01-02 15:04:05"))
	}
	return nil
}

func resolveRCFile() (string, error) {
	shell := os.Getenv("SHELL")
	if shell == "" {
		return "", fmt.Errorf("$SHELL not set")
	}
	home, _ := os.UserHomeDir()
	switch filepath.Base(shell) {
	case "zsh":
		return filepath.Join(home, ".zshrc"), nil
	case "bash":
		rc := filepath.Join(home, ".bashrc")
		if _, err := os.Stat(rc); err != nil {
			return filepath.Join(home, ".bash_profile"), nil
		}
		return rc, nil
	case "fish":
		return filepath.Join(home, ".config", "fish", "config.fish"), nil
	}
	return filepath.Join(home, "."+filepath.Base(shell)+"rc"), nil
}

func viewConfig() error {
	rcFile, err := resolveRCFile()
	if err != nil {
		return err
	}

	pager := os.Getenv("PAGER")
	if pager == "" {
		if _, e := exec.LookPath("bat"); e == nil {
			pager = "bat"
		} else {
			pager = "less"
		}
	}

	var batArgs []string
	if filepath.Base(pager) == "bat" {
		batArgs = []string{"--language", "zsh", "--style", "full", rcFile}
	} else {
		batArgs = []string{rcFile}
	}

	c := exec.Command(pager, batArgs...)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	c.Stdin = os.Stdin
	return c.Run()
}

func backupConfig() error {
	rcFile, err := resolveRCFile()
	if err != nil {
		return err
	}

	home, _ := os.UserHomeDir()
	bdir := filepath.Join(home, backupDir)
	if err := os.MkdirAll(bdir, 0755); err != nil {
		return err
	}

	ts := time.Now().Format("20060102_150405")
	dest := filepath.Join(bdir, filepath.Base(rcFile)+"."+ts)
	if err := copyFile(rcFile, dest); err != nil {
		return fmt.Errorf("backing up: %w", err)
	}
	fmt.Printf("  ✓ Backed up to: %s\n", dest)
	return nil
}

func restoreConfig(id string) error {
	home, _ := os.UserHomeDir()
	bdir := filepath.Join(home, backupDir)

	entries, err := os.ReadDir(bdir)
	if err != nil {
		return fmt.Errorf("no backups found at %s", bdir)
	}

	// Sort newest first
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name() > entries[j].Name()
	})

	var target string
	if id == "" {
		// Use most recent
		if len(entries) == 0 {
			return fmt.Errorf("no backups available")
		}
		target = filepath.Join(bdir, entries[0].Name())
		fmt.Printf("  Restoring most recent backup: %s\n", entries[0].Name())
	} else {
		for _, e := range entries {
			if strings.Contains(e.Name(), id) {
				target = filepath.Join(bdir, e.Name())
				break
			}
		}
		if target == "" {
			return fmt.Errorf("backup %q not found", id)
		}
	}

	rcFile, err := resolveRCFile()
	if err != nil {
		return err
	}

	// Backup current before restoring
	ts := time.Now().Format("20060102_150405")
	preDest := filepath.Join(bdir, filepath.Base(rcFile)+".pre-restore."+ts)
	_ = copyFile(rcFile, preDest)
	fmt.Printf("  Current config saved to: %s\n", preDest)

	if err := copyFile(target, rcFile); err != nil {
		return fmt.Errorf("restoring: %w", err)
	}
	fmt.Printf("  ✓ Restored: %s → %s\n", target, rcFile)
	fmt.Printf("  Reload shell: source ~/.zshrc\n")
	return nil
}

func listBackups() error {
	home, _ := os.UserHomeDir()
	bdir := filepath.Join(home, backupDir)

	entries, err := os.ReadDir(bdir)
	if err != nil {
		fmt.Printf("  (no backups yet — create one with: shellconfig backup)\n")
		return nil
	}

	if len(entries) == 0 {
		fmt.Println("  (no backups yet)")
		return nil
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name() > entries[j].Name()
	})

	fmt.Printf("\n  Backups — %s\n\n", bdir)
	for _, e := range entries {
		info, _ := e.Info()
		fmt.Printf("  %-50s  %d bytes\n", e.Name(), info.Size())
	}
	fmt.Println()
	return nil
}

func diffConfigs(a, b string) error {
	home, _ := os.UserHomeDir()
	rcFile, _ := resolveRCFile()

	resolveArg := func(arg string) string {
		if arg == "" {
			return rcFile
		}
		if filepath.IsAbs(arg) {
			return arg
		}
		bdir := filepath.Join(home, backupDir)
		entries, _ := os.ReadDir(bdir)
		for _, e := range entries {
			if strings.Contains(e.Name(), arg) {
				return filepath.Join(bdir, e.Name())
			}
		}
		return arg
	}

	fileA := resolveArg(a)
	fileB := resolveArg(b)

	differ := "diff"
	if _, err := exec.LookPath("delta"); err == nil {
		differ = "delta"
	}

	c := exec.Command(differ, fileA, fileB)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	c.Run() // diff exits 1 on differences — not an error
	return nil
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
