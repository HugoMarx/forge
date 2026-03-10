package helper

import "charm.land/bubbles/v2/key"


type KeyMap struct {
	Up key.Binding
	Down key.Binding
	Help key.Binding
	Quit key.Binding
	Enter key.Binding
	S key.Binding
	R key.Binding
}

var Keys = KeyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "Move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "Move down"),
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
		key.WithHelp("enter", "Open project"),
	),
	S:  key.NewBinding(
		key.WithKeys("s"),
		key.WithHelp("s", "Show Docker infos"),
	),
	R:  key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "Run Docker env"),
	),
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Enter, k.S, k.R, k.Help, k.Quit}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Enter, k.S}, // first column
		{k.Help, k.Quit},                // second column
	}
}
