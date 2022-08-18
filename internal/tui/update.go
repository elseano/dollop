package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cast"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		statusBarHeight := lipgloss.Height(m.statusView())
		height := m.height - statusBarHeight

		listViewWidth := cast.ToInt(80 * float64(m.width))
		listWidth := listViewWidth - listViewStyle.GetHorizontalFrameSize()
		m.list.SetSize(listWidth, height)

		detailViewWidth := m.width - listViewWidth
		m.viewport = viewport.New(detailViewWidth, height)
		m.viewport.MouseWheelEnabled = true
		m.viewport.SetContent(m.viewportContent(m.viewport.Width))
	case tea.KeyMsg:
		cmds = append(cmds, m.handleKeys(msg))
	case scanMsg:
		m.list.SetItems(msg.lines)
		m.viewport.SetContent(m.viewportContent(m.viewport.Width))

	default:
		m.statusLine = fmt.Sprintf("Msg: %+v", msg)
	}

	cmds = append(cmds, m.scanLogs())
	return m, tea.Batch(cmds...)
}

func (m *Model) handleKeys(msg tea.KeyMsg) tea.Cmd {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg.Type {

	case tea.KeyCtrlC:
		cmd = tea.Quit
		cmds = append(cmds, cmd)
	case tea.KeyUp, tea.KeyDown, tea.KeyLeft, tea.KeyRight:
		m.list, cmd = m.list.Update(msg)
		cmds = append(cmds, cmd)
		m.statusLine = fmt.Sprintf("Cmd: %#v", cmds)
		m.viewport.GotoTop()
		m.viewport.SetContent(m.viewportContent(m.viewport.Width))
	}

	return tea.Batch(cmds...)
}
