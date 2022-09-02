package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

type logGroup struct {
	title        string
	description  string
	groupValue   string
	timestamp    time.Time
	lines        []logLine
	selectedLine int
}

var faintColor = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "#aaaaaa", Dark: "#333333"})

func (i logGroup) Title() string { return i.title }
func (i logGroup) Description() string {
	b := strings.Builder{}
	tally := i.TallyLevels()

	b.WriteString(i.description)
	b.WriteString(" ")

	logCounts := strings.Builder{}

	for _, level := range []string{"trace", "debug", "info", "warning", "error", "fatal"} {

		style := getMessageColor(level)

		if tally[level] > 0 {
			logCounts.WriteString(lipgloss.NewStyle().Foreground(style).MarginLeft(1).Render(fmt.Sprintf("%d", tally[level])))
		}
	}

	if logCounts.Len() > 0 {
		b.WriteString(faintColor.Render(" - "))
		b.WriteString(logCounts.String())
	}

	return b.String()
}
func (i logGroup) FilterValue() string { return "" }

func (i logGroup) TallyLevels() map[string]int {
	result := map[string]int{}

	for _, line := range i.lines {
		cv := result[line.level]
		cv++
		result[line.level] = cv
	}

	return result
}
