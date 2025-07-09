package tui

import "github.com/charmbracelet/bubbles/key"

// based on https://github.com/dlvhdr/gh-dash/blob/main/ui/keys/keys.go

// keys that the user can press
type KeyMap struct {
	// navigation keys (leave the navigation behaviour)
	// Up          key.Binding
	// Down        key.Binding
	// PageDown    key.Binding
	// PageUp      key.Binding
	// NextSection key.Binding
	// PrevSection key.Binding
	// help & quit keys (leave the default help & quit behaviour)
	// Help key.Binding
	// Quit key.Binding
	// custom app key
	SelectItem   key.Binding
	NextItem     key.Binding
	PreviousItem key.Binding
}

// bindings for each navigation key
//func (k KeyMap) NavigationKeys() []key.Binding {
//	return []key.Binding{
//		k.Up,
//		k.Down,
//		k.PrevSection,
//		k.NextSection,
//		k.PageDown,
//		k.PageUp,
//	}
//}

// bindings for quit and help keys (leave the default help & quit behaviour)
//unc (k KeyMap) QuitAndHelpKeys() []key.Binding {
//	return []key.Binding{k.Help, k.Quit}
//

// bindings for each custom app key
func (k KeyMap) AppKeys() []key.Binding {
	return []key.Binding{
		k.SelectItem,
		k.NextItem,
		k.PreviousItem,
	}
}

var Keys = KeyMap{
	// navigation keys
	//Up: key.NewBinding(
	//	key.WithKeys("up"),
	//	key.WithHelp("↑", "up"),
	//),
	//Down: key.NewBinding(
	//	key.WithKeys("down"),
	//	key.WithHelp("↓", "down"),
	//),
	//PrevSection: key.NewBinding(
	//	key.WithKeys("left"),
	//	key.WithHelp("←", "previous section"),
	//),
	//NextSection: key.NewBinding(
	//	key.WithKeys("right"),
	//	key.WithHelp("→", "next section"),
	//),
	//PageDown: key.NewBinding(
	//	key.WithKeys("ctrl+d"),
	//	key.WithHelp("Ctrl+d", "page down"),
	//),
	//PageUp: key.NewBinding(
	//	key.WithKeys("ctrl+u"),
	//	key.WithHelp("Ctrl+u", "page up"),
	//),
	// help & quit keys (leave the default help & quit behaviour)
	//Help: key.NewBinding(
	//	key.WithKeys("?"),
	//	key.WithHelp("?", "help"),
	//),
	//Quit: key.NewBinding(
	//	key.WithKeys("q", "esc", "ctrl+c"),
	//	key.WithHelp("q", "quit"),
	//),
	// custom app keys
	SelectItem: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select kubeconfig"),
	),
	NextItem: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "Next item"),
	),
	PreviousItem: key.NewBinding(
		key.WithKeys("shift+tab"),
		key.WithHelp("shift+tab", "Previous item"),
	),
}

// leave the default help
//func (k KeyMap) ShortHelp() []key.Binding {
//	return []key.Binding{k.Help}
//}

// override the Full Help (when you press ?)
func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		// k.NavigationKeys(),  // leave default navigation behaviour
		k.AppKeys(),
		// k.QuitAndHelpKeys(), // leave default help & quit behaviour
	}
}
