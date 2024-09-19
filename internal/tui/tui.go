package tui

import (
	"strconv"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	defaultWidth  = 30
	defaultHeight = 30
)

var appStyle = lipgloss.NewStyle().Padding(1, 2)

type item struct {
	title   string // title showed in list
	context string // context obtained from the kubeconfig
}

func NewItem(title, context string) item {
	return item{
		title:   title,
		context: context,
	}
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.context }
func (i item) FilterValue() string { return i.title } // field used to filter

// Model defines the model to show
// in this case, it has a list of items, a selected item (choice) and a flag to quit
type Model struct {
	list   list.Model // list of items
	keys   KeyMap     // custom keys to accept
	Choice string     // selected item
}

func NewModel(items []list.Item) *Model {
	// setup list
	l := list.New(items, newListDelegate(Keys), defaultWidth, defaultHeight)
	l.Title = "⎈ Select kubeconfig:"
	l.Styles.Title = lipgloss.NewStyle().Foreground(lipgloss.Color("#038cfc")).Padding(0, 1).Bold(true)
	l.SetFilteringEnabled(true) // enable filtering using Item.FilterValue()

	return &Model{
		list: l,
		keys: Keys,
	}
}

// Init method will be called immediately when the program starts,
// it will do some initialization work, and return a Cmd to tell bubbletea what command to execute
func (m Model) Init() tea.Cmd {
	return nil
}

// Update method is used to respond to external events
// modifies the model, and returns the modified model and a command that bubbletea will execute
// then, bubbletea will call update again, and use the result of the latest tea.Cmd (the tea.Msg) as parameter,
// and will exec the loop
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	// needed so the framework can render the list
	case tea.WindowSizeMsg:
		h, v := appStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)

	// if the user inputs some key
	case tea.KeyMsg:
		// Don't match any of the keys below if we're actively filtering.
		if m.list.FilterState() == list.Filtering {
			break
		}
		// Filter already applied
		if m.list.FilterState() == list.FilterApplied {
			switch {
			case key.Matches(msg, m.keys.SelectItem):
				i, ok := m.list.SelectedItem().(item)
				if ok {
					m.Choice = i.title
				}
				return m, tea.Quit

			case key.Matches(msg, m.keys.NextItem):
				if m.list.Index() == len(m.list.VisibleItems())-1 {
					m.list.ResetSelected()
				} else {
					m.list.CursorDown()
				}
			}
			break
		}

		switch {
		case key.Matches(msg, m.keys.SelectItem):
			i, ok := m.list.SelectedItem().(item)
			if ok {
				m.Choice = i.title
			}
			return m, tea.Quit

		case key.Matches(msg, m.keys.NextItem):
			if m.list.Index() == len(m.list.Items())-1 {
				m.list.ResetSelected()
			} else {
				m.list.CursorDown()
			}
		}
	}

	// update the model
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)

	return m, cmd
}

// View method renders the model into the console
// in this case, calls the View method of the list inside my custom model
func (m Model) View() string {
	index := "index" + strconv.Itoa(m.list.Index())

	len := "len" + strconv.Itoa(len(m.list.VisibleItems()))
	return appStyle.Render(m.list.View() + index + len)
}
