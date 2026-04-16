// Package sheme is the Shell Theme Engine for termbox.
//
// It converts Wondertone palettes — either built-in, or loaded from .wtone
// files — into bash/zsh-sourceable .theme files.
//
// A generated .theme file does two things:
//  1. Exports THEME_* shell variables used by zsh-syntax-highlighting,
//     ZSH_AUTOSUGGEST_HIGHLIGHT_STYLE, and general colour references.
//  2. Emits OSC 4/10/11/12 escape sequences that repaint the terminal
//     emulator's ANSI palette, foreground, background, and cursor.
//
// The ANSI colour mapping mirrors experiments/sheme/generator.go:
//
//	ANSI 0  = Background tone ("Base")
//	ANSI 1  = Red    hue 14°
//	ANSI 2  = Green  hue 142°
//	ANSI 3  = Yellow hue 38°
//	ANSI 4  = Blue   hue 240°
//	ANSI 5  = Magenta hue 320°
//	ANSI 6  = Cyan   hue 196°
//	ANSI 7  = Foreground tone ("Text")
//	ANSI 8-15 = Bright variants (Lighten(10).Saturate(5))
package sheme

import (
	"fmt"
	"strings"

	"github.com/leraniode/wondertone/core"
	"github.com/leraniode/wondertone/palette"
	"github.com/leraniode/wondertone/wtone"
)

// Theme holds the generated content of a .theme file.
type Theme struct {
	Name    string
	Content string // full bash/zsh source text
}

// FromPalette generates a Theme from a *palette.Palette.
// The theme name is derived from the palette name (lowercased, spaces→underscores).
func FromPalette(p *palette.Palette) *Theme {
	name := strings.ToLower(strings.ReplaceAll(p.Name(), " ", "_"))
	return &Theme{
		Name:    name,
		Content: generate(name, p),
	}
}

// FromWToneFile loads a .wtone file from disk and generates a Theme from it.
func FromWToneFile(path string) (*Theme, error) {
	p, err := wtone.LoadWTone(path)
	if err != nil {
		return nil, fmt.Errorf("loading .wtone file %q: %w", path, err)
	}
	return FromPalette(p), nil
}

// FromWToneBytes parses .wtone content from bytes (e.g. go:embed) and generates a Theme.
func FromWToneBytes(data []byte) (*Theme, error) {
	p, err := wtone.ParseWTone(data)
	if err != nil {
		return nil, fmt.Errorf("parsing .wtone bytes: %w", err)
	}
	return FromPalette(p), nil
}

// Scaffold returns a blank .theme template with placeholder Catppuccin Mocha
// values, ready for manual editing.
func Scaffold(name string) string {
	return fmt.Sprintf(`# sheme theme: %s
# Scaffolded by termbox — edit hex values then apply with:
#   termbox sheme apply %s

export THEME_NAME=%q

# ── Core Palette ─────────────────────────────────────────────────────────────
export THEME_BG="#1e1e2e"
export THEME_FG="#cdd6f4"
export THEME_PRIMARY="#cba6f7"
export THEME_SECONDARY="#89b4fa"
export THEME_ACCENT="#cba6f7"
export THEME_SUCCESS="#a6e3a1"
export THEME_WARNING="#f9e2af"
export THEME_ERROR="#f38ba8"
export THEME_MUTED="#585b70"
export THEME_BORDER="#313244"

# ── ANSI 16 ──────────────────────────────────────────────────────────────────
export THEME_BLACK="#45475a"
export THEME_RED="#f38ba8"
export THEME_GREEN="#a6e3a1"
export THEME_YELLOW="#f9e2af"
export THEME_BLUE="#89b4fa"
export THEME_MAGENTA="#cba6f7"
export THEME_CYAN="#89dceb"
export THEME_WHITE="#cdd6f4"
export THEME_BRIGHT_BLACK="#585b70"
export THEME_BRIGHT_RED="#f38ba8"
export THEME_BRIGHT_GREEN="#a6e3a1"
export THEME_BRIGHT_YELLOW="#f9e2af"
export THEME_BRIGHT_BLUE="#89b4fa"
export THEME_BRIGHT_MAGENTA="#cba6f7"
export THEME_BRIGHT_CYAN="#89dceb"
export THEME_BRIGHT_WHITE="#a6adc8"

# ── ZSH Syntax Highlighting ───────────────────────────────────────────────────
typeset -A ZSH_HIGHLIGHT_STYLES 2>/dev/null || true
ZSH_HIGHLIGHT_STYLES[command]="fg=$THEME_GREEN"
ZSH_HIGHLIGHT_STYLES[builtin]="fg=$THEME_BLUE"
ZSH_HIGHLIGHT_STYLES[function]="fg=$THEME_MAGENTA"
ZSH_HIGHLIGHT_STYLES[alias]="fg=$THEME_MAGENTA,bold"
ZSH_HIGHLIGHT_STYLES[string]="fg=$THEME_YELLOW"
ZSH_HIGHLIGHT_STYLES[comment]="fg=$THEME_MUTED"
ZSH_HIGHLIGHT_STYLES[path]="fg=$THEME_FG,underline"
ZSH_HIGHLIGHT_STYLES[globbing]="fg=$THEME_WARNING"
ZSH_HIGHLIGHT_STYLES[unknown-token]="fg=$THEME_ERROR,underline"

# ── ZSH Autosuggestions ───────────────────────────────────────────────────────
export ZSH_AUTOSUGGEST_HIGHLIGHT_STYLE="fg=$THEME_MUTED"

# ── OSC Sequences ─────────────────────────────────────────────────────────────
printf "\033]10;$THEME_FG\007"
printf "\033]11;$THEME_BG\007"
printf "\033]12;$THEME_PRIMARY\007"
printf "\033]4;0;$THEME_BLACK\007"
printf "\033]4;1;$THEME_RED\007"
printf "\033]4;2;$THEME_GREEN\007"
printf "\033]4;3;$THEME_YELLOW\007"
printf "\033]4;4;$THEME_BLUE\007"
printf "\033]4;5;$THEME_MAGENTA\007"
printf "\033]4;6;$THEME_CYAN\007"
printf "\033]4;7;$THEME_WHITE\007"
printf "\033]4;8;$THEME_BRIGHT_BLACK\007"
printf "\033]4;9;$THEME_BRIGHT_RED\007"
printf "\033]4;10;$THEME_BRIGHT_GREEN\007"
printf "\033]4;11;$THEME_BRIGHT_YELLOW\007"
printf "\033]4;12;$THEME_BRIGHT_BLUE\007"
printf "\033]4;13;$THEME_BRIGHT_MAGENTA\007"
printf "\033]4;14;$THEME_BRIGHT_CYAN\007"
printf "\033]4;15;$THEME_BRIGHT_WHITE\007"
`, name, name, name)
}

// ── Internal generation ───────────────────────────────────────────────────────

// getTone looks up a tone by name from a palette, trying "PaletteName Name"
// first then bare "Name". Returns a zero Tone (Hex() == "") if not found.
func getTone(p *palette.Palette, name string) core.Tone {
	if t, ok := p.Get(p.Name() + " " + name); ok {
		return t
	}
	if t, ok := p.Get(name); ok {
		return t
	}
	return core.Tone{}
}

// safeHex returns t.Hex() if non-empty, otherwise the fallback string.
func safeHex(t core.Tone, fallback string) string {
	if h := t.Hex(); h != "" {
		return h
	}
	return fallback
}

// generate builds the full .theme file content from a palette.
// This is the core of the sheme experiment, pulled into a proper package.
func generate(name string, p *palette.Palette) string {
	bg     := getTone(p, "Base")
	fg     := getTone(p, "Text")
	accent := getTone(p, "Accent")
	cursor := accent

	// Derive ANSI lightness/vibrancy from accent (fallback to sensible defaults)
	l, v := 60.0, 70.0
	if (accent != core.Tone{}) {
		l = accent.Light()
		v = accent.Vibrancy()
	}

	// ANSI hues: Red Green Yellow Blue Magenta Cyan
	ansiHues := []float64{14, 142, 38, 240, 320, 196}
	ansi := make([]core.Tone, 16)

	ansi[0] = bg
	for i, hue := range ansiHues {
		ansi[i+1] = core.New(core.Light(l), core.Vibrancy(v), core.Hue(hue))
	}
	ansi[7]  = fg
	ansi[8]  = bg.Lighten(15)
	for i := 1; i <= 6; i++ {
		ansi[i+8] = ansi[i].Lighten(10).Saturate(5)
	}
	ansi[15] = fg.Lighten(10)

	bgHex  := safeHex(bg,     "#1e1e2e")
	fgHex  := safeHex(fg,     "#cdd6f4")
	priHex := safeHex(accent, "#cba6f7")
	curHex := safeHex(cursor, priHex)

	successHex := safeHex(ansi[2], "#a6e3a1")
	warnHex    := safeHex(ansi[3], "#f9e2af")
	errorHex   := safeHex(ansi[1], "#f38ba8")
	secHex     := safeHex(ansi[4], "#89b4fa")
	mutedHex   := safeHex(bg.Lighten(20), "#585b70")
	borderHex  := safeHex(bg.Lighten(10), "#313244")

	ansiNames := []string{
		"THEME_BLACK", "THEME_RED", "THEME_GREEN", "THEME_YELLOW",
		"THEME_BLUE", "THEME_MAGENTA", "THEME_CYAN", "THEME_WHITE",
		"THEME_BRIGHT_BLACK", "THEME_BRIGHT_RED", "THEME_BRIGHT_GREEN", "THEME_BRIGHT_YELLOW",
		"THEME_BRIGHT_BLUE", "THEME_BRIGHT_MAGENTA", "THEME_BRIGHT_CYAN", "THEME_BRIGHT_WHITE",
	}

	var sb strings.Builder

	fmt.Fprintf(&sb, "# sheme theme: %s\n", name)
	fmt.Fprintf(&sb, "# Generated from Wondertone palette: %s\n", p.Name())
	fmt.Fprintf(&sb, "# Apply: termbox sheme apply %s\n\n", name)

	fmt.Fprintf(&sb, "export THEME_NAME=%q\n\n", name)

	fmt.Fprintf(&sb, "# ── Core Palette ─────────────────────────────────────────────────────────────\n")
	fmt.Fprintf(&sb, "export THEME_BG=%q\n", bgHex)
	fmt.Fprintf(&sb, "export THEME_FG=%q\n", fgHex)
	fmt.Fprintf(&sb, "export THEME_PRIMARY=%q\n", priHex)
	fmt.Fprintf(&sb, "export THEME_SECONDARY=%q\n", secHex)
	fmt.Fprintf(&sb, "export THEME_ACCENT=%q\n", priHex)
	fmt.Fprintf(&sb, "export THEME_SUCCESS=%q\n", successHex)
	fmt.Fprintf(&sb, "export THEME_WARNING=%q\n", warnHex)
	fmt.Fprintf(&sb, "export THEME_ERROR=%q\n", errorHex)
	fmt.Fprintf(&sb, "export THEME_MUTED=%q\n", mutedHex)
	fmt.Fprintf(&sb, "export THEME_BORDER=%q\n\n", borderHex)

	fmt.Fprintf(&sb, "# ── ANSI 16 ──────────────────────────────────────────────────────────────────\n")
	for i, varName := range ansiNames {
		fmt.Fprintf(&sb, "export %s=%q\n", varName, safeHex(ansi[i], "#888888"))
	}

	fmt.Fprintf(&sb, "\n# ── ZSH Syntax Highlighting ───────────────────────────────────────────────────\n")
	fmt.Fprintf(&sb, "typeset -A ZSH_HIGHLIGHT_STYLES 2>/dev/null || true\n")
	fmt.Fprintf(&sb, "ZSH_HIGHLIGHT_STYLES[command]=\"fg=$THEME_GREEN\"\n")
	fmt.Fprintf(&sb, "ZSH_HIGHLIGHT_STYLES[builtin]=\"fg=$THEME_BLUE\"\n")
	fmt.Fprintf(&sb, "ZSH_HIGHLIGHT_STYLES[function]=\"fg=$THEME_MAGENTA\"\n")
	fmt.Fprintf(&sb, "ZSH_HIGHLIGHT_STYLES[alias]=\"fg=$THEME_MAGENTA,bold\"\n")
	fmt.Fprintf(&sb, "ZSH_HIGHLIGHT_STYLES[string]=\"fg=$THEME_YELLOW\"\n")
	fmt.Fprintf(&sb, "ZSH_HIGHLIGHT_STYLES[comment]=\"fg=$THEME_MUTED\"\n")
	fmt.Fprintf(&sb, "ZSH_HIGHLIGHT_STYLES[path]=\"fg=$THEME_FG,underline\"\n")
	fmt.Fprintf(&sb, "ZSH_HIGHLIGHT_STYLES[globbing]=\"fg=$THEME_WARNING\"\n")
	fmt.Fprintf(&sb, "ZSH_HIGHLIGHT_STYLES[unknown-token]=\"fg=$THEME_ERROR,underline\"\n")

	fmt.Fprintf(&sb, "\n# ── ZSH Autosuggestions ───────────────────────────────────────────────────────\n")
	fmt.Fprintf(&sb, "export ZSH_AUTOSUGGEST_HIGHLIGHT_STYLE=\"fg=$THEME_MUTED\"\n")

	fmt.Fprintf(&sb, "\n# ── OSC Sequences ────────────────────────────────────────────────────────────\n")
	fmt.Fprintf(&sb, "printf \"\\033]10;%s\\007\"\n", fgHex)
	fmt.Fprintf(&sb, "printf \"\\033]11;%s\\007\"\n", bgHex)
	fmt.Fprintf(&sb, "printf \"\\033]12;%s\\007\"\n", curHex)
	for i := 0; i < 16; i++ {
		fmt.Fprintf(&sb, "printf \"\\033]4;%d;%s\\007\"\n", i, safeHex(ansi[i], "#888888"))
	}

	return sb.String()
}
