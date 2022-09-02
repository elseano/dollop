package tui

import (
	"fmt"
	"sort"
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

func (m Model) logsView() string {
	return m.logs.View()
}

func (m Model) statusView() string {
	return m.statusLine + " " + m.Help.View(m.keyMap)
}

func (m Model) detailView() string {
	return m.detail.View()
}

func (m Model) View() string {
	if m.focusLog == nil {
		return lipgloss.JoinVertical(lipgloss.Left,
			lipgloss.JoinHorizontal(lipgloss.Top, m.listView(), m.logsView()),
			m.statusView(),
		)
	} else {
		return lipgloss.JoinVertical(lipgloss.Left,
			lipgloss.JoinHorizontal(lipgloss.Top, m.listView(), m.detailView()),
			m.statusView(),
		)
	}
}

var (
	brightColor = lipgloss.AdaptiveColor{Light: "#000000", Dark: "#FFFFFF"}
	normalColor = lipgloss.AdaptiveColor{Light: "#333333", Dark: "#CCCCCC"}
	dimColor    = lipgloss.AdaptiveColor{Light: "#666666", Dark: "#666666"}

	infoColor    = lipgloss.AdaptiveColor{Light: "#5db6d7", Dark: "#2983a3"}
	warningColor = lipgloss.AdaptiveColor{Light: "#f0bd32", Dark: "#80610a"}
	errorColor   = lipgloss.AdaptiveColor{Light: "#fd4b98", Dark: "#900934"}
	fatalColor   = lipgloss.AdaptiveColor{Light: "#fd4b4b", Dark: "#900909"}
)

var messageColors = map[string]lipgloss.AdaptiveColor{
	"erro":    errorColor,
	"fata":    fatalColor,
	"warn":    warningColor,
	"info":    infoColor,
	"default": normalColor,
}

func getMessageColor(level string) lipgloss.AdaptiveColor {
	if len(level) > 4 {
		if col, ok := messageColors[level[0:4]]; ok {
			return col
		}
	}

	if col, ok := messageColors[level]; ok {
		return col
	}

	return messageColors["default"]
}

var keyStyle = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "#43BF6D", Dark: "#73F59F"}).Width(25)
var dataStyle = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"})

var tagNameStyle = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "#4e51b7", Dark: "#484cb0"})
var tagSoloStyle = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "#5db6d7", Dark: "#5db6d7"})
var tagValueStyle = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "#999999", Dark: "#999999"})

func (m Model) detailContent(width int) string {
	var builder strings.Builder
	if it := m.logs.SelectedItem(); it != nil {
		line := it.(logLine)

		builder.WriteString(line.String(false))
		builder.WriteString("\n\n")
		builder.WriteString(renderMetadata(line.data, 0))
	}

	return wordwrap.String(builder.String(), width)
}

func trimString(s string, length int) string {
	if len(s) > length {
		return s[0:length]
	}
	return s
}

func renderMetadata(s interface{}, indentLevel int) string {
	switch s := s.(type) {
	case string:
		return s

	case int:
		return fmt.Sprintf("%d", s)

	case float64:
		return fmt.Sprintf("%f", s)

	case []string:
		return strings.Join(s, ", ")

	case []int:
		b := []string{}
		for _, i := range s {
			b = append(b, fmt.Sprintf("%d", i))
		}

		return strings.Join(b, ", ")

	case []float64:
		b := []string{}
		for _, i := range s {
			b = append(b, fmt.Sprintf("%f", i))
		}

		return strings.Join(b, ", ")

	case map[string]interface{}:
		keys := []string{}
		for k, _ := range s {
			keys = append(keys, k)
		}

		sort.Strings(keys)

		builder := strings.Builder{}

		indentWidth := 8 * indentLevel
		indentStr := strings.Repeat(" ", indentWidth)

		if indentWidth > 0 {
			// Move to next line before writing key values, as we've already written the key of the parent.
			builder.WriteString("\n")
		}

		for _, k := range keys {
			builder.WriteString(indentStr)
			builder.WriteString(keyStyle.Render(k))
			builder.WriteString(dataStyle.Render(renderMetadata(s[k], indentLevel+1)))
			builder.WriteString("\n")
		}

		return builder.String()

	case []interface{}:
		b := []string{}
		for _, i := range s {
			b = append(b, fmt.Sprintf("%v", i))
		}

		return strings.Join(b, ", ")

	}

	return fmt.Sprintf("%T", s)
}
