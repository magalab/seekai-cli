package display

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type pagerModel struct {
	content string
	ready   bool
	vp      viewport.Model
}

func RunViewport(text string) error {
	p := tea.NewProgram(pagerModel{content: text})
	_, err := p.Run()
	return err
}

func (m pagerModel) Init() tea.Cmd {
	return nil
}

func (m pagerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		headerHeight := 1
		footerHeight := 1
		if !m.ready {
			m.vp = viewport.New(msg.Width, msg.Height-headerHeight-footerHeight)
			m.vp.SetContent(m.content)
			m.ready = true
		} else {
			m.vp.Width = msg.Width
			m.vp.Height = msg.Height - headerHeight - footerHeight
		}
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit
		}
	}
	var cmd tea.Cmd
	m.vp, cmd = m.vp.Update(msg)
	return m, cmd
}

func (m pagerModel) View() string {
	if !m.ready {
		return m.content
	}
	footer := lipgloss.NewStyle().Foreground(lipgloss.Color("8")).Render("up/down scroll, q quit")
	return fmt.Sprintf("%s\n%s", m.vp.View(), padRight(footer, m.vp.Width))
}

func padRight(value string, width int) string {
	if len(value) >= width {
		return value
	}
	return value + strings.Repeat(" ", width-len(value))
}
