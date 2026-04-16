package main

// sheme command — Shell Theme Engine
//
// Thin CLI wrapper around internal/sheme.
// Supports:
//   - Built-in Wondertone palettes  (termbox sheme from <palette-name>)
//   - .wtone files passed directly  (termbox sheme from-wtone <file.wtone>)
//   - .wtone files in assets/wtone/ (termbox sheme from-wtone <name>)
//   - Manual scaffolding            (termbox sheme new <n>)
//   - Applying and listing          (termbox sheme apply / list / info)

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/dominionthedev/termbox/internal/envutil"
	"github.com/dominionthedev/termbox/internal/registry"
	isheme "github.com/dominionthedev/termbox/internal/sheme"
	"github.com/leraniode/wondertone/palette/builtin"
	"github.com/spf13/cobra"
)

var shemeCmd = &cobra.Command{
	Use:   "sheme [list|palettes|from|from-wtone|new|apply|info]",
	Short: "Shell theme engine — create and apply .theme files",
	Long: `sheme turns Wondertone palettes into bash/zsh .theme files.

  termbox sheme list                list installed themes
  termbox sheme palettes            list built-in Wondertone palettes
  termbox sheme from <palette>      generate a theme from a built-in palette
  termbox sheme from-wtone <file>   generate a theme from a .wtone file
                                    (path, or name inside assets/wtone/)
  termbox sheme new <n>          scaffold a blank theme for manual editing
  termbox sheme apply <n>        set the active theme
  termbox sheme info <n>         show colour slots in a theme`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		action := args[0]
		home, err := registry.FindHome()
		if err != nil {
			return err
		}
		themeDir := filepath.Join(home, "assets", "themes")

		switch action {
		case "list":
			return shemeList(themeDir, home)
		case "palettes":
			return shemePalettes()
		case "from":
			if len(args) < 2 {
				return fmt.Errorf("usage: termbox sheme from <palette-name>")
			}
			return shemeFrom(themeDir, strings.Join(args[1:], " "))
		case "from-wtone":
			if len(args) < 2 {
				return fmt.Errorf("usage: termbox sheme from-wtone <file.wtone>")
			}
			return shemeFromWTone(home, themeDir, args[1])
		case "new":
			if len(args) < 2 {
				return fmt.Errorf("usage: termbox sheme new <n>")
			}
			return shemeNew(themeDir, args[1])
		case "apply":
			if len(args) < 2 {
				return fmt.Errorf("usage: termbox sheme apply <n>")
			}
			return shemeApply(home, themeDir, args[1])
		case "info":
			if len(args) < 2 {
				return fmt.Errorf("usage: termbox sheme info <n>")
			}
			return shemeInfo(themeDir, args[1])
		default:
			return fmt.Errorf("unknown action %q — use: list, palettes, from, from-wtone, new, apply, info", action)
		}
	},
}

// ── palette commands ──────────────────────────────────────────────────────────

func shemePalettes() error {
	fmt.Println("\n  Built-in Wondertone palettes:\n")
	for _, p := range builtin.All() {
		fmt.Printf("  %-24s  %s\n", p.Name(), p.Description())
	}
	fmt.Println()
	return nil
}

func shemeList(themeDir, home string) error {
	entries, err := os.ReadDir(themeDir)
	if err != nil {
		fmt.Println("  (no themes yet — create one with: termbox sheme from <palette>)")
		return nil
	}
	active := shemeActiveTheme(home)
	fmt.Printf("\n  Themes — %s\n\n", themeDir)
	found := false
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".theme") {
			name := strings.TrimSuffix(e.Name(), ".theme")
			marker := "  "
			if name == active {
				marker = "▶ "
			}
			fmt.Printf("  %s%s\n", marker, name)
			found = true
		}
	}
	if !found {
		fmt.Println("  (no .theme files yet)")
	}
	fmt.Println()
	return nil
}

// ── generation commands ───────────────────────────────────────────────────────

func shemeFrom(themeDir, paletteName string) error {
	if err := os.MkdirAll(themeDir, 0755); err != nil {
		return fmt.Errorf("creating themes dir: %w", err)
	}

	var found interface {
		Name() string
		Description() string
	}
	for _, p := range builtin.All() {
		if strings.EqualFold(p.Name(), paletteName) {
			found = p
			break
		}
	}
	if found == nil {
		return fmt.Errorf("palette %q not found — run 'termbox sheme palettes'", paletteName)
	}

	// Use the actual palette.Palette type for generation
	for _, p := range builtin.All() {
		if strings.EqualFold(p.Name(), paletteName) {
			theme := isheme.FromPalette(p)
			return writeTheme(themeDir, theme)
		}
	}
	return nil
}

func shemeFromWTone(home, themeDir, arg string) error {
	if err := os.MkdirAll(themeDir, 0755); err != nil {
		return fmt.Errorf("creating themes dir: %w", err)
	}

	// Resolve the .wtone file: try as direct path first, then look in assets/wtone/
	wtoneFile := resolveWToneFile(home, arg)
	if wtoneFile == "" {
		return fmt.Errorf(".wtone file not found: %q\nLooked in: %s and %s",
			arg,
			arg,
			filepath.Join(home, "assets", "wtone"),
		)
	}

	theme, err := isheme.FromWToneFile(wtoneFile)
	if err != nil {
		return err
	}
	return writeTheme(themeDir, theme)
}

// resolveWToneFile returns the absolute path to a .wtone file.
// It accepts:
//  1. An absolute or relative path to a .wtone file
//  2. A bare name (with or without .wtone extension) that lives in assets/wtone/
func resolveWToneFile(home, arg string) string {
	// Try as a direct file path first
	if filepath.IsAbs(arg) {
		if _, err := os.Stat(arg); err == nil {
			return arg
		}
	}
	// Try relative path from cwd
	if _, err := os.Stat(arg); err == nil {
		abs, _ := filepath.Abs(arg)
		return abs
	}
	// Try with .wtone extension
	withExt := arg
	if !strings.HasSuffix(arg, ".wtone") {
		withExt = arg + ".wtone"
	}
	if _, err := os.Stat(withExt); err == nil {
		abs, _ := filepath.Abs(withExt)
		return abs
	}

	// Try inside assets/wtone/
	wtoneDir := filepath.Join(home, "assets", "wtone")
	candidates := []string{
		filepath.Join(wtoneDir, arg),
		filepath.Join(wtoneDir, withExt),
	}
	for _, c := range candidates {
		if _, err := os.Stat(c); err == nil {
			return c
		}
	}
	return ""
}

func shemeNew(themeDir, name string) error {
	if err := os.MkdirAll(themeDir, 0755); err != nil {
		return fmt.Errorf("creating themes dir: %w", err)
	}
	path := filepath.Join(themeDir, name+".theme")
	if _, err := os.Stat(path); err == nil {
		return fmt.Errorf("theme %q already exists at %s", name, path)
	}
	if err := os.WriteFile(path, []byte(isheme.Scaffold(name)), 0644); err != nil {
		return fmt.Errorf("writing theme: %w", err)
	}
	fmt.Printf("  ✓ Scaffolded: %s\n  Edit hex values, then: termbox sheme apply %s\n", path, name)
	return nil
}

// ── apply / info ──────────────────────────────────────────────────────────────

func shemeApply(home, themeDir, name string) error {
	themePath := filepath.Join(themeDir, name+".theme")
	if _, err := os.Stat(themePath); err != nil {
		return fmt.Errorf("theme %q not found — run 'termbox sheme list'", name)
	}

	// Write theme.env — sourced by the shell at startup
	themeEnv := filepath.Join(home, "config", "theme.env")
	content := fmt.Sprintf(
		"# theme.env — managed by 'termbox sheme apply'\nexport SHELL_THEME=%q\n[[ -f %q ]] && source %q\n",
		name, themePath, themePath,
	)
	if err := os.WriteFile(themeEnv, []byte(content), 0644); err != nil {
		return fmt.Errorf("writing theme.env: %w", err)
	}

	// Also reflect in termbox.env
	envPath := filepath.Join(home, "config", "termbox.env")
	if _, err := os.Stat(envPath); err == nil {
		_ = envutil.UpdateVar(envPath, "SHELL_THEME", name)
	}

	fmt.Printf("  ✓ Active theme: %s\n  Reload shell: source ~/.zshrc\n", name)
	return nil
}

func shemeInfo(themeDir, name string) error {
	path := filepath.Join(themeDir, name+".theme")
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("theme %q not found — run 'termbox sheme list'", name)
	}
	fmt.Printf("\n  %s.theme\n\n", name)
	for _, line := range strings.Split(string(data), "\n") {
		t := strings.TrimSpace(line)
		if strings.HasPrefix(t, "export ") || strings.HasPrefix(t, "printf ") {
			fmt.Printf("  %s\n", line)
		}
	}
	fmt.Println()
	return nil
}

// ── helpers ───────────────────────────────────────────────────────────────────

func shemeActiveTheme(home string) string {
	data, _ := os.ReadFile(filepath.Join(home, "config", "theme.env"))
	for _, line := range strings.Split(string(data), "\n") {
		if strings.HasPrefix(line, "export SHELL_THEME=") {
			return strings.Trim(strings.TrimPrefix(line, "export SHELL_THEME="), `"'`)
		}
	}
	return os.Getenv("SHELL_THEME")
}

func writeTheme(themeDir string, t *isheme.Theme) error {
	path := filepath.Join(themeDir, t.Name+".theme")
	if err := os.WriteFile(path, []byte(t.Content), 0644); err != nil {
		return fmt.Errorf("writing theme: %w", err)
	}
	fmt.Printf("  ✓ Generated: %s\n  Apply with: termbox sheme apply %s\n", path, t.Name)
	return nil
}

func init() {
	rootCmd.AddCommand(shemeCmd)
}
