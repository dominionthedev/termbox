package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/dominionthedev/termbox/internal/registry"
	"github.com/dominionthedev/termbox/internal/settings"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var configCmd = &cobra.Command{
	Use:   "config [list|get|set|edit]",
	Short: "Read and write termbox settings (config/settings.yaml)",
	Long: `Manage termbox settings. This is NOT the component registry.

  termbox config list              print all current settings
  termbox config get <key>         read a setting  e.g. theme.active
  termbox config set <key> <value> write a setting e.g. termbox config set notes.folder ~/notes
  termbox config edit              open settings.yaml in $EDITOR`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		action := args[0]
		home, err := registry.FindHome()
		if err != nil {
			return err
		}
		switch action {
		case "list":
			return configList(home)
		case "get":
			if len(args) < 2 {
				return fmt.Errorf("usage: termbox config get <key>")
			}
			return configGet(home, args[1])
		case "set":
			if len(args) < 3 {
				return fmt.Errorf("usage: termbox config set <key> <value>")
			}
			return configSet(home, args[1], strings.Join(args[2:], " "))
		case "edit":
			editor := os.Getenv("EDITOR")
			if editor == "" {
				editor = "nvim"
			}
			c := exec.Command(editor, filepath.Join(home, "config", "settings.yaml"))
			c.Stdout = os.Stdout
			c.Stderr = os.Stderr
			c.Stdin = os.Stdin
			return c.Run()
		default:
			return fmt.Errorf("unknown action %q — use: list, get, set, edit", action)
		}
	},
}

func configList(home string) error {
	s, err := settings.Load(home)
	if err != nil {
		return fmt.Errorf("loading settings: %w", err)
	}
	data, _ := yaml.Marshal(s)
	fmt.Println()
	for _, line := range strings.Split(strings.TrimSpace(string(data)), "\n") {
		fmt.Printf("  %s\n", line)
	}
	fmt.Println()
	return nil
}

func configGet(home, key string) error {
	s, err := settings.Load(home)
	if err != nil {
		return err
	}
	data, _ := yaml.Marshal(s)
	var m map[string]interface{}
	_ = yaml.Unmarshal(data, &m)

	val := dotGet(m, strings.Split(key, "."))
	if val == nil {
		return fmt.Errorf("key %q not found — run 'termbox config list'", key)
	}
	fmt.Printf("  %s = %v\n", key, val)
	return nil
}

func dotGet(m map[string]interface{}, parts []string) interface{} {
	if len(parts) == 0 || m == nil {
		return nil
	}
	v, ok := m[parts[0]]
	if !ok {
		return nil
	}
	if len(parts) == 1 {
		return v
	}
	next, _ := v.(map[string]interface{})
	return dotGet(next, parts[1:])
}

func configSet(home, key, value string) error {
	s, err := settings.Load(home)
	if err != nil {
		return err
	}

	boolVal := value == "true" || value == "1" || value == "yes"

	switch key {
	case "display.color":     s.Display.Color = boolVal
	case "display.icons":     s.Display.Icons = boolVal
	case "display.pager":     s.Display.Pager = boolVal
	case "banner.enabled":    s.Banner.Enabled = boolVal
	case "banner.file":       s.Banner.File = value
	case "banner.color":      s.Banner.Color = boolVal
	case "notes.folder":      s.Notes.Folder = value
	case "notes.extension":   s.Notes.Extension = value
	case "notes.editor":      s.Notes.Editor = value
	case "theme.active":      s.Theme.Active = value
	case "theme.auto_apply":  s.Theme.AutoApply = boolVal
	case "powerup.active":    s.Powerup.Active = value
	case "powerup.auto_detect": s.Powerup.AutoDetect = boolVal
	case "build.output_dir":  s.Build.OutputDir = value
	case "build.go_flags":    s.Build.GoFlags = value
	default:
		return fmt.Errorf("unknown key %q\nRun 'termbox config list' to see all keys", key)
	}

	if err := settings.Save(s, home); err != nil {
		return err
	}
	fmt.Printf("  ✓ %s = %s\n", key, value)
	return nil
}

func init() {
	rootCmd.AddCommand(configCmd)
}
