package tui

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
)

// customize the behaviour when the Update cycle (user input) happens
// we create a default delegate and override some functions, like help,
// to show some keys there...
// we can override more behaviour
func newListDelegate(keys KeyMap) list.DefaultDelegate {
	d := list.NewDefaultDelegate()

	// keys to show in help menu (additionally to the default)
	// to override this, we would have to reimplement the help.KeyMap interface (not sure)
	help := []key.Binding{keys.SelectItem}

	d.ShortHelpFunc = func() []key.Binding {
		return help
	}

	d.FullHelpFunc = func() [][]key.Binding {
		return [][]key.Binding{help}
	}

	return d
}
