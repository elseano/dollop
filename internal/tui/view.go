package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/wordwrap"
)

var (
	listViewStyle = lipgloss.NewStyle().
			PaddingRight(1).
			MarginRight(1).
			Border(lipgloss.RoundedBorder(), false, true, false, false)

	dividerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#9B9B9B", Dark: "#5C5C5C"})
)

func (m Model) listView() string {
	return listViewStyle.Render(m.list.View())
}

func (m Model) detailView() string {
	return m.viewport.View()
}

func (m Model) statusView() string {
	return m.statusLine
}

func (m Model) View() string {
	return lipgloss.JoinVertical(lipgloss.Left,
		lipgloss.JoinHorizontal(lipgloss.Top, m.listView(), m.detailView()),
		m.statusView(),
	)
}

func (m Model) viewportContent(width int) string {
	var builder strings.Builder
	if it := m.list.SelectedItem(); it != nil {
		for _, line := range it.(*item).lines {
			builder.WriteString(line.message + "\n")
		}
	} else {
		builder.WriteString("No item selected")
	}

	return wordwrap.String(builder.String(), width)
}
