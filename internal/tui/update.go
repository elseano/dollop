package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cast"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmds []tea.Cmd
		cmd  tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		statusBarHeight := lipgloss.Height(m.statusView())
		height := m.height - statusBarHeight

		listViewWidth := cast.ToInt(0.30 * float64(m.width))
		if listViewWidth < 30 {
			listViewWidth = 30
		}

		listWidth := listViewWidth - listViewStyle.GetHorizontalFrameSize()
		m.list.SetSize(listWidth, height)

		m.rightSideWidth = m.width - m.list.Width()

		if m.focusLog == nil {
			m.logs.SetWidth(m.rightSideWidth)
			m.logs.SetHeight(height - 1)
		} else {
			m.detail = viewport.New(m.rightSideWidth, height)
			m.detail.MouseWheelEnabled = true
			m.detail.SetContent(m.detailContent(m.detail.Width))
		}

		lines, selLine := m.generateLogItems()
		m.logs.SetItems(lines)
		m.logs.Select(selLine)

	case tea.MouseMsg:
		switch m.focus {
		case "groups":
			m.list, cmd = m.list.Update(msg)
			cmds = append(cmds, cmd)
		case "logs":
			if m.focusLog != nil {
				m.detail, cmd = m.detail.Update(msg)
				cmds = append(cmds, cmd)
			} else {
				m.logs, cmd = m.logs.Update(msg)
				cmds = append(cmds, cmd)
			}
		}

	case tea.KeyMsg:
		cmds = append(cmds, m.handleKeys(msg))

	case disconnectedMsg:
		m.disconnected = true
		m.SetStatus("Process has terminated")

	case scanMsg:
		sel, selected := m.list.SelectedItem().(*logGroup)

		m.list.SetItems(msg.lines)

		if selected {
			for index, item := range msg.lines {
				if sel.title == item.(*logGroup).title {
					m.list.Select(index)
					break
				}
			}
		}

		lines, selLine := m.generateLogItems()
		m.logs.SetItems(lines)
		m.logs.Select(selLine)

		if m.focus == "groups" {
			m.setKeysForIndex(&m.list)
		} else if m.focus == "logs" && m.focusLog == nil {
			m.setKeysForIndex(&m.logs)
		} else if m.focus == "" {
			m.focusOnGroups()
		}

		if msg.status != "" {
			m.SetStatus(msg.status)
		} else if m.statusLine == "" {
			m.SetStatus("Logs receiving")
		}

	default:
		m.statusLine = fmt.Sprintf("Msg: %+v", msg)
	}

	if !m.disconnected {
		cmds = append(cmds, m.scanLogs())
	}

	return m, tea.Batch(cmds...)
}

func (m *Model) handleKeys(msg tea.KeyMsg) tea.Cmd {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch {

	case key.Matches(msg, m.keyMap.Quit):
		cmd = tea.Quit
		cmds = append(cmds, cmd)

	case key.Matches(msg, m.keyMap.GoToStart):
		switch m.focus {
		case "groups":
			m.list.Select(0)
		case "logs":
			if m.focusLog == nil {
				m.logs.Select(0)
			}
		}

	case key.Matches(msg, m.keyMap.GoToEnd):
		switch m.focus {
		case "groups":
			m.list.Select(len(m.list.Items()) - 1)
		case "logs":
			if m.focusLog == nil {
				m.logs.Select(len(m.logs.Items()) - 1)
			}
		}

	case key.Matches(msg, m.keyMap.CursorDown, m.keyMap.CursorUp, m.keyMap.NextPage, m.keyMap.PrevPage):
		switch m.focus {
		case "groups":
			m.list, cmd = m.list.Update(msg)
			cmds = append(cmds, cmd)
			items, selected := m.generateLogItems()
			m.logs.SetItems(items)
			m.logs.Select(selected)

			m.setKeysForIndex(&m.list)
		case "logs":
			if m.focusLog != nil {
				m.detail, cmd = m.detail.Update(msg)
				cmds = append(cmds, cmd)
			} else {
				m.logs, cmd = m.logs.Update(msg)
				cmds = append(cmds, cmd)

				if group, ok := m.list.SelectedItem().(*logGroup); ok {
					group.selectedLine = m.logs.Index()
				}

				m.setKeysForIndex(&m.logs)
			}

		}

	case key.Matches(msg, m.keyMap.Select):
		switch m.focus {
		case "logs":
			log, ok := m.logs.SelectedItem().(logLine)
			if ok {
				m.focusOnLogItem(&log)
			}

		case "groups":
			m.focusOnLogs()

		}

	case key.Matches(msg, m.keyMap.Escape):
		switch m.focus {
		case "logs":
			if m.focusLog != nil {
				m.focusOnLogs()
			} else {
				m.focusOnGroups()
			}
		}

	default:
		// m.statusLine = fmt.Sprintf("%#v", msg)

		// m.viewport.GotoTop()
		// m.viewport.SetContent(m.viewportContent(m.viewport.Width))
	}

	return tea.Batch(cmds...)
}

func (m *Model) setKeysForIndex(l *list.Model) {
	if l.Index() == 0 {
		m.keyMap.CursorUp.SetEnabled(false)
		m.keyMap.PrevPage.SetEnabled(false)
		m.keyMap.CursorDown.SetEnabled(true)
		m.keyMap.NextPage.SetEnabled(true)
	} else if l.Index() == len(l.Items())-1 {
		m.keyMap.CursorUp.SetEnabled(true)
		m.keyMap.PrevPage.SetEnabled(true)
		m.keyMap.CursorDown.SetEnabled(false)
		m.keyMap.NextPage.SetEnabled(false)
	} else {
		m.keyMap.CursorUp.SetEnabled(true)
		m.keyMap.PrevPage.SetEnabled(true)
		m.keyMap.CursorDown.SetEnabled(true)
		m.keyMap.NextPage.SetEnabled(true)
	}
}

func (m *Model) focusOnGroups() {
	m.focus = "groups"
	m.logs.Title = "Logs"
	m.list.Title = "Groups (active)"

	m.keyMap.Escape.SetEnabled(false)
	m.setKeysForIndex(&m.list)
	m.logs.SetDelegate(NewLogLineDelegate(false))
}

func (m *Model) focusOnLogs() {
	m.focus = "logs"
	m.logs.Title = "Logs (active)"
	m.list.Title = "Groups"
	m.focusLog = nil

	m.keyMap.Escape.SetEnabled(true)
	m.keyMap.Select.SetEnabled(true)
	m.setKeysForIndex(&m.logs)
	m.logs.SetDelegate(NewLogLineDelegate(true))
}

func (m *Model) focusOnLogItem(log *logLine) {
	m.focusLog = log

	m.detail = viewport.New(m.rightSideWidth, m.height-1)
	m.detail.MouseWheelEnabled = true
	m.detail.SetContent(m.detailContent(m.detail.Width))

	m.keyMap.Select.SetEnabled(false)
}

func (m Model) generateLogItems() ([]list.Item, int) {
	result := []list.Item{}

	if it, ok := m.list.SelectedItem().(*logGroup); ok && it != nil {
		for _, line := range it.lines {
			result = append(result, line)
		}

		return result, it.selectedLine
	}

	return result, 0
}

func (m *Model) SetStatus(status string) {
	m.statusLine = status
	m.Help.Width = m.width - (len(m.statusLine) + 3)
}
