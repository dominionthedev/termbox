package ui

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/dominionthedev/termbox/internal/styles"
	
	"github.com/charmbracelet/lipgloss"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
)

// SystemPanel renders system stats with history
type SystemPanel struct {
	CPUPercent  float64
	MemPercent  float64
	DiskPercent float64
	CPUHistory  []float64
	MemHistory  []float64
	Pulse       *styles.PulseState
}

func NewSystemPanel() *SystemPanel {
	return &SystemPanel{
		CPUHistory: make([]float64, 0, 30),
		MemHistory: make([]float64, 0, 30),
		Pulse:      styles.NewPulseState(1.0),
	}
}

func (p *SystemPanel) Update() error {
	// CPU
	cpuPercentages, err := cpu.Percent(time.Millisecond*100, false)
	if err == nil && len(cpuPercentages) > 0 {
		p.CPUPercent = cpuPercentages[0]
		p.CPUHistory = append(p.CPUHistory, p.CPUPercent)
		if len(p.CPUHistory) > 30 {
			p.CPUHistory = p.CPUHistory[1:]
		}
	}
	
	// Memory
	memInfo, err := mem.VirtualMemory()
	if err == nil {
		p.MemPercent = memInfo.UsedPercent
		p.MemHistory = append(p.MemHistory, p.MemPercent)
		if len(p.MemHistory) > 30 {
			p.MemHistory = p.MemHistory[1:]
		}
	}
	
	// Disk
	diskInfo, err := disk.Usage("/")
	if err == nil {
		p.DiskPercent = diskInfo.UsedPercent
	}
	
	return nil
}

func (p *SystemPanel) Render(width int) string {
	var b strings.Builder
	
	// Header
	b.WriteString(styles.GetHeaderStyle().Render(styles.IconCPU + " SYSTEM METRICS"))
	b.WriteString("\n\n")
	
	// CPU with sparkline
	cpuBar := styles.RenderProgressBar(p.CPUPercent, 25, p.Pulse)
	b.WriteString(styles.RenderMetric("CPU", fmt.Sprintf("%5.1f%%", p.CPUPercent)))
	b.WriteString("\n")
	b.WriteString("  " + cpuBar)
	b.WriteString("\n")
	if len(p.CPUHistory) > 0 {
		sparkline := styles.RenderSparkline(p.CPUHistory, 25)
		b.WriteString("  " + sparkline)
	}
	b.WriteString("\n\n")
	
	// Memory with sparkline
	memBar := styles.RenderProgressBar(p.MemPercent, 25, nil)
	b.WriteString(styles.RenderMetric("Memory", fmt.Sprintf("%5.1f%%", p.MemPercent)))
	b.WriteString("\n")
	b.WriteString("  " + memBar)
	b.WriteString("\n")
	if len(p.MemHistory) > 0 {
		sparkline := styles.RenderSparkline(p.MemHistory, 25)
		b.WriteString("  " + sparkline)
	}
	b.WriteString("\n\n")
	
	// Disk
	diskBar := styles.RenderProgressBar(p.DiskPercent, 25, nil)
	b.WriteString(styles.RenderMetric("Disk", fmt.Sprintf("%5.1f%%", p.DiskPercent)))
	b.WriteString("\n")
	b.WriteString("  " + diskBar)
	b.WriteString("\n\n")
	
	// System info
	cores, _ := cpu.Counts(false)
	b.WriteString(lipgloss.NewStyle().Foreground(styles.CurrentTheme.TextSecondary).Render(
		fmt.Sprintf("OS: %s | Cores: %d", runtime.GOOS, cores),
	))
	
	return styles.GetPanelStyle().Width(width - 4).Render(b.String())
}

// NetworkPanel renders network stats
type NetworkPanel struct {
	BytesSent    uint64
	BytesRecv    uint64
	PacketsSent  uint64
	PacketsRecv  uint64
	LastUpdate   time.Time
	SendRate     float64
	RecvRate     float64
	SendHistory  []float64
	RecvHistory  []float64
}

func NewNetworkPanel() *NetworkPanel {
	return &NetworkPanel{
		LastUpdate:  time.Now(),
		SendHistory: make([]float64, 0, 30),
		RecvHistory: make([]float64, 0, 30),
	}
}

func (p *NetworkPanel) Update() error {
	netIO, err := net.IOCounters(false)
	if err != nil || len(netIO) == 0 {
		return err
	}
	
	now := time.Now()
	elapsed := now.Sub(p.LastUpdate).Seconds()
	
	if elapsed > 0 && p.BytesSent > 0 {
		p.SendRate = float64(netIO[0].BytesSent-p.BytesSent) / elapsed
		p.RecvRate = float64(netIO[0].BytesRecv-p.BytesRecv) / elapsed
		
		p.SendHistory = append(p.SendHistory, p.SendRate/1024/1024)
		p.RecvHistory = append(p.RecvHistory, p.RecvRate/1024/1024)
		
		if len(p.SendHistory) > 30 {
			p.SendHistory = p.SendHistory[1:]
		}
		if len(p.RecvHistory) > 30 {
			p.RecvHistory = p.RecvHistory[1:]
		}
	}
	
	p.BytesSent = netIO[0].BytesSent
	p.BytesRecv = netIO[0].BytesRecv
	p.PacketsSent = netIO[0].PacketsSent
	p.PacketsRecv = netIO[0].PacketsRecv
	p.LastUpdate = now
	
	return nil
}

func (p *NetworkPanel) Render(width int) string {
	var b strings.Builder
	
	b.WriteString(styles.GetHeaderStyle().Render(styles.IconNetwork + " NETWORK"))
	b.WriteString("\n\n")
	
	// Upload
	uploadMB := float64(p.BytesSent) / 1024 / 1024
	b.WriteString(styles.RenderMetric("Upload", fmt.Sprintf("%.1f MB", uploadMB)))
	b.WriteString("\n")
	b.WriteString(styles.RenderMetric("Rate", fmt.Sprintf("%.1f KB/s", p.SendRate/1024)))
	b.WriteString("\n")
	if len(p.SendHistory) > 0 {
		sparkline := styles.RenderSparkline(p.SendHistory, 25)
		b.WriteString("  " + sparkline)
	}
	b.WriteString("\n\n")
	
	// Download
	downloadMB := float64(p.BytesRecv) / 1024 / 1024
	b.WriteString(styles.RenderMetric("Download", fmt.Sprintf("%.1f MB", downloadMB)))
	b.WriteString("\n")
	b.WriteString(styles.RenderMetric("Rate", fmt.Sprintf("%.1f KB/s", p.RecvRate/1024)))
	b.WriteString("\n")
	if len(p.RecvHistory) > 0 {
		sparkline := styles.RenderSparkline(p.RecvHistory, 25)
		b.WriteString("  " + sparkline)
	}
	b.WriteString("\n\n")
	
	// Packets
	b.WriteString(lipgloss.NewStyle().Foreground(styles.CurrentTheme.TextSecondary).Render(
		fmt.Sprintf("Packets: %s %d | %s %d", styles.IconUp, p.PacketsSent, styles.IconDown, p.PacketsRecv),
	))
	
	return styles.GetPanelStyle().Width(width - 4).Render(b.String())
}

// GitPanel renders detailed git information
type GitPanel struct {
	Branch       string
	Dirty        bool
	Ahead        int
	Behind       int
	Modified     int
	Staged       int
	Untracked    int
	LastCommit   string
	RepoPath     string
}

func NewGitPanel() *GitPanel {
	return &GitPanel{}
}

func (p *GitPanel) Update() error {
	// Check if in git repo
	cmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
	if err := cmd.Run(); err != nil {
		return err
	}
	
	// Get current directory
	cmd = exec.Command("git", "rev-parse", "--show-toplevel")
	if output, err := cmd.Output(); err == nil {
		p.RepoPath = strings.TrimSpace(string(output))
	}
	
	// Get branch
	cmd = exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	if output, err := cmd.Output(); err == nil {
		p.Branch = strings.TrimSpace(string(output))
	}
	
	// Get ahead/behind
	cmd = exec.Command("git", "rev-list", "--left-right", "--count", "HEAD...@{upstream}")
	if output, err := cmd.Output(); err == nil {
		parts := strings.Fields(string(output))
		if len(parts) == 2 {
			fmt.Sscanf(parts[0], "%d", &p.Ahead)
			fmt.Sscanf(parts[1], "%d", &p.Behind)
		}
	}
	
	// Get file status
	cmd = exec.Command("git", "status", "--porcelain")
	if output, err := cmd.Output(); err == nil {
		p.Modified = 0
		p.Staged = 0
		p.Untracked = 0
		
		lines := strings.Split(string(output), "\n")
		for _, line := range lines {
			if len(line) < 2 {
				continue
			}
			if line[:2] == "??" {
				p.Untracked++
			} else if line[0] != ' ' && line[0] != '?' {
				p.Staged++
			} else if line[1] != ' ' {
				p.Modified++
			}
		}
		
		p.Dirty = len(output) > 0
	}
	
	// Get last commit
	cmd = exec.Command("git", "log", "-1", "--pretty=%s")
	if output, err := cmd.Output(); err == nil {
		p.LastCommit = strings.TrimSpace(string(output))
		if len(p.LastCommit) > 40 {
			p.LastCommit = p.LastCommit[:37] + "..."
		}
	}
	
	return nil
}

func (p *GitPanel) Render(width int) string {
	var b strings.Builder
	
	b.WriteString(styles.GetHeaderStyle().Render(styles.IconGit + " GIT STATUS"))
	b.WriteString("\n\n")
	
	// Branch with status
	branchStyle := lipgloss.NewStyle().Foreground(styles.CurrentTheme.Primary).Bold(true)
	b.WriteString(styles.IconGitBranch + " ")
	b.WriteString(branchStyle.Render(p.Branch))
	
	if p.Dirty {
		b.WriteString(" " + styles.RenderBadge("DIRTY", "warning"))
	} else {
		b.WriteString(" " + styles.RenderBadge("CLEAN", "success"))
	}
	b.WriteString("\n\n")
	
	// Ahead/Behind
	if p.Ahead > 0 || p.Behind > 0 {
		syncStatus := ""
		if p.Ahead > 0 {
			syncStatus += fmt.Sprintf("%s %d ahead", styles.IconUp, p.Ahead)
		}
		if p.Behind > 0 {
			if p.Ahead > 0 {
				syncStatus += " | "
			}
			syncStatus += fmt.Sprintf("%s %d behind", styles.IconDown, p.Behind)
		}
		b.WriteString(lipgloss.NewStyle().Foreground(styles.CurrentTheme.Warning).Render(syncStatus))
		b.WriteString("\n\n")
	}
	
	// File counts
	if p.Modified > 0 {
		b.WriteString(fmt.Sprintf("  %s %d modified\n", styles.IconWarning, p.Modified))
	}
	if p.Staged > 0 {
		b.WriteString(fmt.Sprintf("  %s %d staged\n", styles.IconCheck, p.Staged))
	}
	if p.Untracked > 0 {
		b.WriteString(fmt.Sprintf("  ? %d untracked\n", p.Untracked))
	}
	
	if p.Modified == 0 && p.Staged == 0 && p.Untracked == 0 {
		b.WriteString(lipgloss.NewStyle().Foreground(styles.CurrentTheme.Success).Render("  Working tree clean"))
	}
	
	b.WriteString("\n\n")
	
	// Last commit
	if p.LastCommit != "" {
		b.WriteString(lipgloss.NewStyle().Foreground(styles.CurrentTheme.TextSecondary).Render("Last commit:"))
		b.WriteString("\n")
		b.WriteString(lipgloss.NewStyle().Foreground(styles.CurrentTheme.TextPrimary).Italic(true).Render("  " + p.LastCommit))
	}
	
	return styles.GetPanelStyle().Width(width - 4).Render(b.String())
}

// QuickActionsPanel renders quick action buttons
type QuickActionsPanel struct {
	Selected int
	Actions  []Action
}

type Action struct {
	Name        string
	Icon        string
	Description string
	Command     string
}

func NewQuickActionsPanel() *QuickActionsPanel {
	return &QuickActionsPanel{
		Actions: []Action{
			{Name: "New Session", Icon: styles.IconTerminal, Description: "Create tmux session", Command: "bt session new"},
			{Name: "Jump Project", Icon: styles.IconFolder, Description: "Navigate to project", Command: "bt project jump"},
			{Name: "Git Status", Icon: styles.IconGit, Description: "Repository overview", Command: "bt git status"},
			{Name: "Switch Theme", Icon: styles.IconStar, Description: "Change color scheme", Command: "bt theme switch"},
			{Name: "Capture", Icon: styles.IconRocket, Description: "Screenshot terminal", Command: "bt capture"},
			{Name: "Monitor", Icon: styles.IconProcess, Description: "Process monitor", Command: "bt monitor"},
		},
	}
}

func (p *QuickActionsPanel) Render(width int) string {
	var b strings.Builder
	
	b.WriteString(styles.GetHeaderStyle().Render(styles.IconBolt + " QUICK ACTIONS"))
	b.WriteString("\n\n")
	
	for i, action := range p.Actions {
		var style lipgloss.Style
		if i == p.Selected {
			style = styles.GetListItemSelectedStyle()
		} else {
			style = styles.GetListItemStyle()
		}
		
		line := fmt.Sprintf("%s %s", action.Icon, action.Name)
		b.WriteString(style.Render(line))
		b.WriteString("\n")
		
		if i == p.Selected {
			desc := "  " + lipgloss.NewStyle().Foreground(styles.CurrentTheme.TextSecondary).Render(action.Description)
			b.WriteString(desc)
			b.WriteString("\n")
		}
	}
	
	return styles.GetPanelStyle().Width(width - 4).Render(b.String())
}

// StatusPanel renders current status
type StatusPanel struct {
	GitBranch     string
	GitDirty      bool
	TmuxSessions  int
	ActiveProject string
	Timestamp     time.Time
	Uptime        time.Duration
	ActiveTheme   string
}

func NewStatusPanel() *StatusPanel {
	return &StatusPanel{
		GitBranch:     "main",
		GitDirty:      false,
		TmuxSessions:  0,
		ActiveProject: "~",
		Timestamp:     time.Now(),
		ActiveTheme:   styles.CurrentTheme.Name,
	}
}

func (p *StatusPanel) Update() {
	p.Timestamp = time.Now()
	p.ActiveTheme = styles.CurrentTheme.Name
	
	// Get tmux session count
	cmd := exec.Command("tmux", "list-sessions")
	if output, err := cmd.Output(); err == nil {
		p.TmuxSessions = len(strings.Split(string(output), "\n")) - 1
	}
}

func (p *StatusPanel) Render(width int) string {
	var b strings.Builder
	
	b.WriteString(styles.GetHeaderStyle().Render(styles.IconInfo + " STATUS"))
	b.WriteString("\n\n")
	
	// Git status
	gitStatus := styles.IconGit + " " + p.GitBranch
	if p.GitDirty {
		gitStatus += " " + styles.GetStatusWarning().Render("✗")
	} else {
		gitStatus += " " + styles.GetStatusOK().Render("✓")
	}
	b.WriteString(lipgloss.NewStyle().Foreground(styles.CurrentTheme.TextPrimary).Render(gitStatus))
	b.WriteString("\n\n")
	
	// Tmux sessions
	tmuxStatus := fmt.Sprintf("%s %d sessions", styles.IconTerminal, p.TmuxSessions)
	b.WriteString(lipgloss.NewStyle().Foreground(styles.CurrentTheme.TextPrimary).Render(tmuxStatus))
	b.WriteString("\n\n")
	
	// Active project
	projectStatus := fmt.Sprintf("%s %s", styles.IconFolder, p.ActiveProject)
	b.WriteString(lipgloss.NewStyle().Foreground(styles.CurrentTheme.TextPrimary).Render(projectStatus))
	b.WriteString("\n\n")
	
	// Theme
	themeStatus := fmt.Sprintf("%s %s", styles.IconStar, p.ActiveTheme)
	b.WriteString(lipgloss.NewStyle().Foreground(styles.CurrentTheme.Tertiary).Render(themeStatus))
	b.WriteString("\n\n")
	
	// Timestamp
	timeStr := p.Timestamp.Format("15:04:05")
	b.WriteString(lipgloss.NewStyle().Foreground(styles.CurrentTheme.TextSecondary).Render(styles.IconClock + " " + timeStr))
	
	return styles.GetPanelStyle().Width(width - 4).Render(b.String())
}

// WelcomeBanner renders the main welcome banner
func RenderWelcomeBanner() string {
	banner := `
╔══════════════════════════════════════════════╗
║                                              ║
║   ██████╗ ███████╗ █████╗ ██╗   ██╗████████╗║
║   ██╔══██╗██╔════╝██╔══██╗██║   ██║╚══██╔══╝║
║   ██████╔╝█████╗  ███████║██║   ██║   ██║   ║
║   ██╔══██╗██╔══╝  ██╔══██║██║   ██║   ██║   ║
║   ██████╔╝███████╗██║  ██║╚██████╔╝   ██║   ║
║   ╚═════╝ ╚══════╝╚═╝  ╚═╝ ╚═════╝    ╚═╝   ║
║                                              ║
║          T E R M   D A S H B O A R D         ║
║                                              ║
╚══════════════════════════════════════════════╝
`
	
	gradient := styles.GetGradient(banner)
	return lipgloss.NewStyle().MarginBottom(1).Render(gradient)
}

// HelpBar renders the help/keybindings at the bottom
func RenderHelpBar(width int) string {
	helps := []string{
		"↑/↓ navigate",
		"enter select",
		"tab switch panel",
		"t theme",
		"r refresh",
		"q quit",
	}
	
	helpText := strings.Join(helps, " • ")
	helpStyle := lipgloss.NewStyle().
		Foreground(styles.CurrentTheme.TextDim).
		Italic(true).
		Width(width).
		Align(lipgloss.Center)
	
	return helpStyle.Render(helpText)
}
