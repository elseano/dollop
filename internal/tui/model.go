package tui

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/elseano/dollop/internal/config"
)

type Model struct {
	width, height int

	list           list.Model
	logs           list.Model
	focus          string
	statusLine     string
	rightSideWidth int

	focusLog *logLine
	detail   viewport.Model

	Help   help.Model
	keyMap KeyMap

	disconnected bool

	config config.Config
}

func New(config config.Config) (*Model, error) {
	l := list.New(nil, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Requests"
	l.SetShowHelp(false)
	l.SetShowStatusBar(false)
	l.SetShowFilter(false)
	l.SetFilteringEnabled(false)

	logs := list.New(nil, NewLogLineDelegate(false), 0, 0)
	logs.SetShowHelp(false)
	logs.SetShowStatusBar(false)
	logs.SetShowFilter(false)
	logs.SetFilteringEnabled(false)

	return &Model{
		list:         l,
		logs:         logs,
		config:       config,
		Help:         help.New(),
		keyMap:       DefaultKeyMap(),
		disconnected: false,
	}, nil
}
