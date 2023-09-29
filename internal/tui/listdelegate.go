package tui

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

// customize the behaviour when the Update cycle (user input) happens
// we create a default delegate and override behaviour:
// - show new keys in the help section
// - items style
func newListDelegate(keys KeyMap) list.DefaultDelegate {
	d := list.NewDefaultDelegate()

	// keys to show in help menu (additionally to the default)
	// to override this, we would have to reimplement the help.KeyMap interface (not sure)
	help := []key.Binding{keys.SelectItem}

	// override help
	d.ShortHelpFunc = func() []key.Binding {
		return help
	}

	// override expanded help
	d.FullHelpFunc = func() [][]key.Binding {
		return [][]key.Binding{help}
	}

	// override items style
	d.Styles.NormalTitle = lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#1a1a1a", Dark: "#dddddd"}).
		Padding(0, 0, 0, 2)

	d.Styles.NormalDesc = d.Styles.NormalTitle.Copy().
		Foreground(lipgloss.AdaptiveColor{Light: "#A49FA5", Dark: "#777777"})

	d.Styles.SelectedTitle = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), false, false, false, true).
		BorderForeground(lipgloss.AdaptiveColor{Light: "#000000", Dark: "#038cfc"}).
		Foreground(lipgloss.AdaptiveColor{Light: "#000000", Dark: "#038cfc"}).
		Padding(0, 0, 0, 1)

	d.Styles.SelectedDesc = d.Styles.SelectedTitle.Copy().
		Foreground(lipgloss.AdaptiveColor{Light: "#000000", Dark: "#777777"})

	d.Styles.DimmedTitle = lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#A49FA5", Dark: "#777777"}).
		Padding(0, 0, 0, 2)

	d.Styles.DimmedDesc = d.Styles.DimmedTitle.Copy().
		Foreground(lipgloss.AdaptiveColor{Light: "#C2B8C2", Dark: "#4D4D4D"})

	d.Styles.FilterMatch = lipgloss.NewStyle().Underline(true)

	return d
}
