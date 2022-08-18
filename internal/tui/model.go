package tui

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/elseano/sl/internal/config"
)

type Model struct {
	width, height int

	list       list.Model
	viewport   viewport.Model
	statusLine string

	config config.Config
}

func New(config config.Config) (*Model, error) {
	l := list.New(nil, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Requests"
	l.SetShowHelp(false)
	l.SetShowStatusBar(false)
	l.SetShowFilter(false)
	l.SetFilteringEnabled(false)

	return &Model{
		list:   l,
		config: config,
	}, nil
}
