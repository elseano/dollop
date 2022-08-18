package tui

import "time"

type item struct {
	title       string
	description string
	timestamp   time.Time
	lines       []logLine
}

type logLine struct {
	message string
	data    map[string]interface{}
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.description }
func (i item) FilterValue() string { return "" }
