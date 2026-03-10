package helper

import "charm.land/bubbles/v2/key"

type KeyMap struct {
	Up    key.Binding
	Down  key.Binding
	Help  key.Binding
	Quit  key.Binding
	Enter key.Binding
	s     key.Binding
	u     key.Binding
	U     key.Binding
	d     key.Binding
}

var Keys = KeyMap{
	Up: key.NewBinding(
		key.WithKeys("up"),
		key.WithHelp("↑", "Move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down"),
		key.WithHelp("↓", "Move down"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "Toggle help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q", "Quit"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("Enter", "Open project"),
	),
	s: key.NewBinding(
		key.WithKeys("s"),
		key.WithHelp("s", "Show Docker infos"),
	),
	u: key.NewBinding(
		key.WithKeys("u"),
		key.WithHelp("u", "Run Docker env"),
	),
	U: key.NewBinding(
		key.WithKeys("U"),
		key.WithHelp("U", "Run Docker env and open"),
	),
	d: key.NewBinding(
		key.WithKeys("d"),
		key.WithHelp("d", "Unmount Docker env"),
	),
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Enter, k.s, k.u, k.U, k.d, k.Help, k.Quit}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.s, k.u, k.U},
		{k.d, k.Enter, k.Help, k.Quit},
	}
}
