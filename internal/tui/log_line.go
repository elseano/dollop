package tui

import (
	"io"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type logTag struct {
	name  string
	value string
}

type logLine struct {
	level     string
	timestamp time.Time
	message   string
	data      map[string]interface{}
	tags      []logTag
}

func (line logLine) FilterValue() string {
	return line.message
}

type logLineDelegate struct{ isActive bool }

func NewLogLineDelegate(active bool) logLineDelegate {
	return logLineDelegate{isActive: active}
}

var unselectedColor = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}
var selectedColor = lipgloss.AdaptiveColor{Light: "#EE6FF8", Dark: "#EE6FF8"}
var selectedBackgroundColor = lipgloss.AdaptiveColor{Light: "#f49efa", Dark: "#890792"}
var selectedBackgroundStyle = lipgloss.NewStyle().Background(selectedBackgroundColor)
var unselectedBackgroundStyle = lipgloss.NewStyle()

// func DebugColors() {
// 	fmt.Printf("Info: %#v\n", labelStyles["info"])
// 	fmt.Printf("Selected: %#v\n", selectedBackgroundStyle)
// 	fmt.Printf("Merged: %#v\n", selectedBackgroundStyle.Inherit(labelStyles["info"]))
// 	fmt.Printf("Merged2: %#v\n", labelStyles["info"].Inherit(selectedBackgroundStyle))
// }

func (lineDelegate logLineDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	line := item.(logLine)
	builder := strings.Builder{}

	if index == m.Index() && lineDelegate.isActive {
		builder.WriteString("  " + line.String(true))
	} else {
		builder.WriteString("  " + line.String(false))
	}

	w.Write([]byte(builder.String()))
}

func (line logLine) String(selected bool) string {
	builder := strings.Builder{}
	lineStyle := unselectedBackgroundStyle
	if selected {
		lineStyle = selectedBackgroundStyle
	}

	messageColour := getMessageColor(line.level)
	textColor := normalColor

	if messageColour == messageColors["default"] {
		textColor = dimColor
	}

	messageStyle := lipgloss.NewStyle().Foreground(textColor)
	labelStyle := lipgloss.NewStyle().Background(messageColour).Foreground(textColor)

	if line.level == "debug" || line.level == "trace" {
		labelStyle = lipgloss.NewStyle().Foreground(textColor)
	}

	levelStr := trimString(strings.ToUpper(line.level), 4)
	levelStr = " " + levelStr + " "

	if selected {
		builder.WriteString(labelStyle.Background(lineStyle.GetBackground()).Render(levelStr))
	} else {
		builder.WriteString(labelStyle.Render(levelStr))
	}

	builder.WriteString(lineStyle.Render(strings.Repeat(" ", 9-lipgloss.Width(levelStr))))
	builder.WriteString(lineStyle.Inherit(messageStyle).Render(strings.ReplaceAll(line.message, "\n", " ")))

	if len(line.tags) > 0 {
		builder.WriteString(lineStyle.Render("   "))
	}

	for _, t := range line.tags {
		if t.value != "" {
			builder.WriteString(lineStyle.Render(" "))
			builder.WriteString(lineStyle.Inherit(tagNameStyle).Render(t.name))
			builder.WriteString(lineStyle.Render(" "))
			builder.WriteString(lineStyle.Inherit(tagValueStyle).Render(t.value))
		} else {
			builder.WriteString(lineStyle.Render(" "))
			builder.WriteString(lineStyle.Inherit(tagSoloStyle).Render(t.name))
		}
	}

	return builder.String()
}

func (lineDelegate logLineDelegate) Height() int {
	return 1
}

func (lineDelegate logLineDelegate) Spacing() int {
	return 0
}

func (lineDelegate logLineDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	return nil
}
