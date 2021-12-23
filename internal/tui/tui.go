package tui

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Model defines the model to show
// in this case, it has a list, a selected item (choice) and a flag to quit
type Model struct {
	list list.Model
	// items    []Item
	choice   chan string
	quitting bool
}

func NewModel(l list.Model, choice chan string) *Model {
	return &Model{
		list:   l,
		choice: choice,
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
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	// if the
	case ErrMsg:
		return m, nil
	// if the user inputs some key
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c", "q":
			close(m.choice)
			m.quitting = true
			return m, tea.Quit

		case "enter":
			i, ok := m.list.SelectedItem().(Item)
			if ok {
				m.choice <- string(i)
			}
			return m, tea.Quit
		}
	}
	// handle the default keys
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

// View method renders the model into the console
// in this case, calls the View method of the list inside my custom model
func (m Model) View() string {
	//if m.choice != "" {
	//	return quitTextStyle.Render(fmt.Sprintf("Selecting %s", m.choice))
	//}
	if m.quitting {
		return quitTextStyle.Render("Exiting...")
	}
	return "\n" + m.list.View()
}

const ListHeight = 14

var (
	TitleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("12"))
	PaginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	HelpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)

type ErrMsg struct{ err error }

func NewErrMsg(err error) *ErrMsg {
	return &ErrMsg{
		err: err,
	}
}

type Item string

func (i Item) FilterValue() string { return "" }

type ItemDelegate struct{}

func (d ItemDelegate) Height() int                               { return 1 }
func (d ItemDelegate) Spacing() int                              { return 0 }
func (d ItemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d ItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(Item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s string) string {
			return selectedItemStyle.Render("> " + s)
		}
	}

	fmt.Fprintf(w, fn(str))
}
