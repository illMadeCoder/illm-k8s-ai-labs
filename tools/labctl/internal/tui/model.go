package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/illmadecoder/labctl/internal/k8s"
	"github.com/illmadecoder/labctl/internal/tutorial"
)

// Model is the bubbletea model for the tutorial TUI.
type Model struct {
	tutorial       *tutorial.Tutorial
	experiment     *k8s.ExperimentInfo
	kubeconfigPath string

	currentPage int
	width       int
	height      int

	renderer *glamour.TermRenderer
}

// NewModel creates a new TUI model.
func NewModel(tut *tutorial.Tutorial, exp *k8s.ExperimentInfo, kubeconfigPath string) Model {
	r, _ := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(100),
	)

	return Model{
		tutorial:       tut,
		experiment:     exp,
		kubeconfigPath: kubeconfigPath,
		renderer:       r,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "n", "right", "l", " ":
			if m.currentPage < len(m.tutorial.Pages)-1 {
				m.currentPage++
			}
		case "p", "left", "h":
			if m.currentPage > 0 {
				m.currentPage--
			}
		case "home", "g":
			m.currentPage = 0
		case "end", "G":
			m.currentPage = len(m.tutorial.Pages) - 1
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		if m.width > 0 {
			r, _ := glamour.NewTermRenderer(
				glamour.WithAutoStyle(),
				glamour.WithWordWrap(m.width-4),
			)
			m.renderer = r
		}
	}

	return m, nil
}

func (m Model) View() string {
	if m.width == 0 {
		return "Loading..."
	}

	var sections []string

	// Header
	sections = append(sections, m.renderHeader())

	// Services bar
	if m.experiment != nil && len(m.experiment.Services) > 0 {
		sections = append(sections, m.renderServices())
	}

	// Kubeconfig bar
	if m.kubeconfigPath != "" {
		sections = append(sections, m.renderKubeconfig())
	}

	// Content
	sections = append(sections, m.renderContent())

	// Footer
	sections = append(sections, m.renderFooter())

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

func (m Model) renderHeader() string {
	page := m.tutorial.Pages[m.currentPage]
	title := m.tutorial.Title
	if title == "" {
		title = m.tutorial.Name
	}

	pageInfo := fmt.Sprintf("[%d/%d]", m.currentPage+1, len(m.tutorial.Pages))

	left := titleStyle.Render(title)
	right := pageIndicatorStyle.Render(pageInfo)

	gap := m.width - lipgloss.Width(left) - lipgloss.Width(right)
	if gap < 1 {
		gap = 1
	}

	header := left + strings.Repeat(" ", gap) + right

	// Info line
	var infoParts []string
	if m.experiment != nil {
		infoParts = append(infoParts, fmt.Sprintf("Phase: %s", m.experiment.Phase))
		if len(m.experiment.Targets) > 0 {
			infoParts = append(infoParts, fmt.Sprintf("Cluster: %s", m.experiment.Targets[0].ClusterName))
		}
		if m.experiment.TTLDays > 0 {
			infoParts = append(infoParts, fmt.Sprintf("TTL: %dd", m.experiment.TTLDays))
		}
	}
	if page.Title != "" && page.Title != title {
		infoParts = append(infoParts, fmt.Sprintf("Section: %s", page.Title))
	}

	info := ""
	if len(infoParts) > 0 {
		info = "\n" + infoStyle.Render(strings.Join(infoParts, " | "))
	}

	return header + info
}

func (m Model) renderServices() string {
	var lines []string
	lines = append(lines, infoStyle.Render("SERVICES:"))
	for _, svc := range m.experiment.Services {
		name := serviceNameStyle.Render(svc.Name + ":")
		endpoint := serviceStyle.Render(svc.Endpoint)
		readyMark := ""
		if !svc.Ready {
			readyMark = " (pending)"
		}
		lines = append(lines, fmt.Sprintf("  %s %s%s", name, endpoint, readyMark))
	}
	return strings.Join(lines, "\n")
}

func (m Model) renderKubeconfig() string {
	return kubeconfigStyle.Render(fmt.Sprintf("KUBECONFIG: export KUBECONFIG=%s", m.kubeconfigPath))
}

func (m Model) renderContent() string {
	page := m.tutorial.Pages[m.currentPage]

	// Calculate available height for content
	headerHeight := 3
	servicesHeight := 0
	if m.experiment != nil && len(m.experiment.Services) > 0 {
		servicesHeight = len(m.experiment.Services) + 1
	}
	kubeconfigHeight := 0
	if m.kubeconfigPath != "" {
		kubeconfigHeight = 1
	}
	footerHeight := 1
	overhead := headerHeight + servicesHeight + kubeconfigHeight + footerHeight + 2

	contentHeight := m.height - overhead
	if contentHeight < 5 {
		contentHeight = 5
	}

	// Render markdown with glamour
	rendered := page.Content
	if m.renderer != nil {
		if out, err := m.renderer.Render(page.Content); err == nil {
			rendered = out
		}
	}

	// Trim to fit
	lines := strings.Split(rendered, "\n")
	if len(lines) > contentHeight {
		lines = lines[:contentHeight]
		lines = append(lines, infoStyle.Render("... (more content, use scrolling in future version)"))
	}

	return strings.Join(lines, "\n")
}

func (m Model) renderFooter() string {
	keys := "[n]ext  [p]rev  [q]uit"
	pageNum := fmt.Sprintf("Page %d of %d", m.currentPage+1, len(m.tutorial.Pages))

	gap := m.width - len(keys) - len(pageNum)
	if gap < 1 {
		gap = 1
	}

	return footerStyle.Render(keys + strings.Repeat(" ", gap) + pageNum)
}
