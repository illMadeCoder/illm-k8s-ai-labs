package tui

import "github.com/charmbracelet/lipgloss"

var (
	// Colors
	accent    = lipgloss.Color("#7D56F4")
	subtle    = lipgloss.Color("#383838")
	highlight = lipgloss.Color("#EE6FF8")
	green     = lipgloss.Color("#04B575")
	yellow    = lipgloss.Color("#FFCC00")
	red       = lipgloss.Color("#FF4444")
	white     = lipgloss.Color("#FAFAFA")
	dim       = lipgloss.Color("#666666")

	// Styles
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(white).
			Background(accent).
			Padding(0, 1)

	infoStyle = lipgloss.NewStyle().
			Foreground(dim)

	serviceStyle = lipgloss.NewStyle().
			Foreground(green).
			Bold(true)

	serviceNameStyle = lipgloss.NewStyle().
				Foreground(yellow).
				Width(16)

	kubeconfigStyle = lipgloss.NewStyle().
			Foreground(highlight).
			Bold(true)

	borderStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(subtle)

	footerStyle = lipgloss.NewStyle().
			Foreground(dim)

	pageIndicatorStyle = lipgloss.NewStyle().
				Foreground(accent).
				Bold(true)
)
