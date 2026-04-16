package styles

import (
	"math"
	"time"

	"github.com/charmbracelet/lipgloss"
)

// Theme represents a color theme
type Theme struct {
	Name       string
	Primary    lipgloss.Color
	Secondary  lipgloss.Color
	Tertiary   lipgloss.Color
	Success    lipgloss.Color
	Warning    lipgloss.Color
	Error      lipgloss.Color
	BgDark     lipgloss.Color
	BgDarker   lipgloss.Color
	BgPanel    lipgloss.Color
	BgAccent   lipgloss.Color
	TextPrimary   lipgloss.Color
	TextSecondary lipgloss.Color
	TextDim       lipgloss.Color
}

// Available themes
var (
	CyberpunkTheme = Theme{
		Name:       "Cyberpunk",
		Primary:    lipgloss.Color("#00f0ff"),
		Secondary:  lipgloss.Color("#ff00ff"),
		Tertiary:   lipgloss.Color("#bd93f9"),
		Success:    lipgloss.Color("#50fa7b"),
		Warning:    lipgloss.Color("#f1fa8c"),
		Error:      lipgloss.Color("#ff5555"),
		BgDark:     lipgloss.Color("#0a0e14"),
		BgDarker:   lipgloss.Color("#050810"),
		BgPanel:    lipgloss.Color("#151a23"),
		BgAccent:   lipgloss.Color("#1a1f2e"),
		TextPrimary:   lipgloss.Color("#e6e6e6"),
		TextSecondary: lipgloss.Color("#8b92a8"),
		TextDim:       lipgloss.Color("#5a5f73"),
	}

	DraculaTheme = Theme{
		Name:       "Dracula",
		Primary:    lipgloss.Color("#8be9fd"),
		Secondary:  lipgloss.Color("#ff79c6"),
		Tertiary:   lipgloss.Color("#bd93f9"),
		Success:    lipgloss.Color("#50fa7b"),
		Warning:    lipgloss.Color("#f1fa8c"),
		Error:      lipgloss.Color("#ff5555"),
		BgDark:     lipgloss.Color("#282a36"),
		BgDarker:   lipgloss.Color("#1e1f29"),
		BgPanel:    lipgloss.Color("#44475a"),
		BgAccent:   lipgloss.Color("#6272a4"),
		TextPrimary:   lipgloss.Color("#f8f8f2"),
		TextSecondary: lipgloss.Color("#6272a4"),
		TextDim:       lipgloss.Color("#44475a"),
	}

	TokyoNightTheme = Theme{
		Name:       "Tokyo Night",
		Primary:    lipgloss.Color("#7aa2f7"),
		Secondary:  lipgloss.Color("#bb9af7"),
		Tertiary:   lipgloss.Color("#9ece6a"),
		Success:    lipgloss.Color("#9ece6a"),
		Warning:    lipgloss.Color("#e0af68"),
		Error:      lipgloss.Color("#f7768e"),
		BgDark:     lipgloss.Color("#1a1b26"),
		BgDarker:   lipgloss.Color("#16161e"),
		BgPanel:    lipgloss.Color("#24283b"),
		BgAccent:   lipgloss.Color("#414868"),
		TextPrimary:   lipgloss.Color("#c0caf5"),
		TextSecondary: lipgloss.Color("#565f89"),
		TextDim:       lipgloss.Color("#414868"),
	}
)

// Current active theme
var CurrentTheme = CyberpunkTheme

// Animated elements
type PulseState struct {
	StartTime time.Time
	Frequency float64
}

func NewPulseState(frequency float64) *PulseState {
	return &PulseState{
		StartTime: time.Now(),
		Frequency: frequency,
	}
}

func (p *PulseState) GetIntensity() float64 {
	elapsed := time.Since(p.StartTime).Seconds()
	return (math.Sin(elapsed*p.Frequency*2*math.Pi) + 1) / 2
}

// Style generators using current theme
func GetAppStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Background(CurrentTheme.BgDark).
		Foreground(CurrentTheme.TextPrimary)
}

func GetPanelStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Background(CurrentTheme.BgPanel).
		Foreground(CurrentTheme.TextPrimary).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(CurrentTheme.Primary).
		Padding(1, 2).
		MarginRight(1).
		MarginBottom(1)
}

func GetPanelActiveStyle() lipgloss.Style {
	return GetPanelStyle().
		BorderForeground(CurrentTheme.Secondary).
		BorderStyle(lipgloss.ThickBorder())
}

func GetHeaderStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Bold(true).
		Foreground(CurrentTheme.Primary).
		Background(CurrentTheme.BgDarker).
		Padding(0, 1).
		MarginBottom(1)
}

func GetTitleStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Bold(true).
		Foreground(CurrentTheme.Secondary).
		Background(CurrentTheme.BgDarker).
		Padding(0, 2).
		MarginBottom(1)
}

// Status styles
func GetStatusOK() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(CurrentTheme.Success).Bold(true)
}

func GetStatusWarning() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(CurrentTheme.Warning).Bold(true)
}

func GetStatusError() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(CurrentTheme.Error).Bold(true)
}

func GetStatusInfo() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(CurrentTheme.Primary).Bold(true)
}

// List styles
func GetListItemStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(CurrentTheme.TextPrimary).
		PaddingLeft(2)
}

func GetListItemSelectedStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(CurrentTheme.Primary).
		Background(CurrentTheme.BgAccent).
		Bold(true).
		PaddingLeft(2)
}

// Icons using Nerd Fonts
const (
	IconCPU       = " "
	IconMemory    = " "
	IconDisk      = " "
	IconNetwork   = "󰛳 "
	IconGit       = " "
	IconGitBranch = " "
	IconFolder    = " "
	IconFile      = " "
	IconTerminal  = " "
	IconCheck     = " "
	IconX         = " "
	IconWarning   = " "
	IconInfo      = " "
	IconStar      = " "
	IconClock     = " "
	IconArrowR    = " "
	IconArrowL    = " "
	IconArrowU    = " "
	IconArrowD    = " "
	IconDot       = "●"
	IconCircle    = "○"
	IconBolt      = "⚡"
	IconFire      = "🔥"
	IconRocket    = "🚀"
	IconUp        = "󰁝 "
	IconDown      = "󰁅 "
	IconSync      = "󰑓 "
	IconProcess   = "󰊕 "
)

// RenderProgressBar with animation support
func RenderProgressBar(percent float64, width int, pulse *PulseState) string {
	filled := int(float64(width) * percent / 100)
	if filled > width {
		filled = width
	}
	
	bar := ""
	for i := 0; i < width; i++ {
		if i < filled {
			bar += "█"
		} else {
			bar += "░"
		}
	}
	
	filledStyle := lipgloss.NewStyle().Foreground(CurrentTheme.Primary)
	emptyStyle := lipgloss.NewStyle().Foreground(CurrentTheme.TextDim)
	
	// Add pulse animation
	if pulse != nil {
		intensity := pulse.GetIntensity()
		if intensity > 0.7 {
			filledStyle = filledStyle.Foreground(CurrentTheme.Secondary)
		}
	}
	
	return filledStyle.Render(bar[:filled]) + emptyStyle.Render(bar[filled:])
}

// RenderMetric displays a labeled metric
func RenderMetric(label, value string) string {
	labelStyle := lipgloss.NewStyle().
		Foreground(CurrentTheme.TextSecondary).
		Width(12)
	valueStyle := lipgloss.NewStyle().
		Foreground(CurrentTheme.Primary).
		Bold(true)
	
	return labelStyle.Render(label) + " " + valueStyle.Render(value)
}

// RenderStatus returns a status indicator with icon
func RenderStatus(status string, ok bool) string {
	if ok {
		return GetStatusOK().Render(IconCheck + " " + status)
	}
	return GetStatusError().Render(IconX + " " + status)
}

// RenderDivider creates a horizontal divider
func RenderDivider(width int) string {
	divider := ""
	for i := 0; i < width; i++ {
		divider += "─"
	}
	return lipgloss.NewStyle().
		Foreground(CurrentTheme.TextDim).
		MarginTop(1).
		MarginBottom(1).
		Render(divider)
}

// GetGradient creates a gradient effect for titles
func GetGradient(text string) string {
	colors := []lipgloss.Color{
		CurrentTheme.Primary,
		CurrentTheme.Tertiary,
		CurrentTheme.Secondary,
	}
	output := ""
	
	for i, char := range text {
		color := colors[i%len(colors)]
		output += lipgloss.NewStyle().Foreground(color).Render(string(char))
	}
	
	return output
}

// RenderSparkline creates a mini graph
func RenderSparkline(values []float64, width int) string {
	if len(values) == 0 {
		return ""
	}
	
	// Find max value
	maxVal := 0.0
	for _, v := range values {
		if v > maxVal {
			maxVal = v
		}
	}
	if maxVal == 0 {
		maxVal = 1
	}
	
	// Sparkline characters
	chars := []string{"▁", "▂", "▃", "▄", "▅", "▆", "▇", "█"}
	
	result := ""
	step := len(values) / width
	if step < 1 {
		step = 1
	}
	
	for i := 0; i < len(values); i += step {
		if len(result) >= width {
			break
		}
		normalized := values[i] / maxVal
		idx := int(normalized * float64(len(chars)-1))
		if idx >= len(chars) {
			idx = len(chars) - 1
		}
		
		color := CurrentTheme.Primary
		if normalized > 0.8 {
			color = CurrentTheme.Error
		} else if normalized > 0.6 {
			color = CurrentTheme.Warning
		}
		
		result += lipgloss.NewStyle().Foreground(color).Render(chars[idx])
	}
	
	return result
}

// RenderBox creates a fancy box
func RenderBox(title, content string, width int) string {
	titleStyle := lipgloss.NewStyle().
		Foreground(CurrentTheme.Secondary).
		Bold(true).
		Padding(0, 1)
	
	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(CurrentTheme.Primary).
		Width(width).
		Padding(1)
	
	header := titleStyle.Render(title)
	return boxStyle.Render(header + "\n\n" + content)
}

// RenderBadge creates a badge
func RenderBadge(text string, style string) string {
	var badgeStyle lipgloss.Style
	
	switch style {
	case "success":
		badgeStyle = lipgloss.NewStyle().
			Foreground(CurrentTheme.BgDark).
			Background(CurrentTheme.Success).
			Bold(true).
			Padding(0, 1)
	case "warning":
		badgeStyle = lipgloss.NewStyle().
			Foreground(CurrentTheme.BgDark).
			Background(CurrentTheme.Warning).
			Bold(true).
			Padding(0, 1)
	case "error":
		badgeStyle = lipgloss.NewStyle().
			Foreground(CurrentTheme.BgDark).
			Background(CurrentTheme.Error).
			Bold(true).
			Padding(0, 1)
	default:
		badgeStyle = lipgloss.NewStyle().
			Foreground(CurrentTheme.BgDark).
			Background(CurrentTheme.Primary).
			Bold(true).
			Padding(0, 1)
	}
	
	return badgeStyle.Render(text)
}

// SetTheme changes the current theme
func SetTheme(theme Theme) {
	CurrentTheme = theme
}

// GetAvailableThemes returns all themes
func GetAvailableThemes() []Theme {
	return []Theme{
		CyberpunkTheme,
		DraculaTheme,
		TokyoNightTheme,
	}
}
