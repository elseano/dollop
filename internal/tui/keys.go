package tui

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	CursorUp   key.Binding
	CursorDown key.Binding
	NextPage   key.Binding
	PrevPage   key.Binding
	GoToStart  key.Binding
	GoToEnd    key.Binding

	Select key.Binding
	Escape key.Binding
	Quit   key.Binding
}

func DefaultKeyMap() KeyMap {
	return KeyMap{
		CursorUp: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "up"),
		),
		CursorDown: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "down"),
		),
		PrevPage: key.NewBinding(
			key.WithKeys("h", "pgup", "b", "u"),
			key.WithHelp("h/pgup", "prev page"),
		),
		NextPage: key.NewBinding(
			key.WithKeys("l", "pgdown", "f", "d"),
			key.WithHelp("l/pgdn", "next page"),
		),
		GoToStart: key.NewBinding(
			key.WithKeys("home", "g"),
			key.WithHelp("g/home", "go to start"),
		),
		GoToEnd: key.NewBinding(
			key.WithKeys("end", "G"),
			key.WithHelp("G/end", "go to end"),
		),

		Select: key.NewBinding(
			key.WithKeys("right", "enter"),
			key.WithHelp("→/enter", "select"),
		),
		Escape: key.NewBinding(
			key.WithKeys("left", "esc"),
			key.WithHelp("←/esc", "back"),
		),

		Quit: key.NewBinding(
			key.WithKeys("ctrl+c"),
			key.WithHelp("ctrl+c", "quit"),
		),
	}
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{
			k.CursorDown,
			k.CursorUp,
			k.NextPage,
			k.PrevPage,
			k.GoToEnd,
			k.GoToStart,
			k.Select,
			k.Escape,
			k.Quit,
		},
	}
}

func (k KeyMap) ShortHelp() []key.Binding {
	return k.FullHelp()[0]
}
