package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ebarped/kubeswap/pkg/kv"
	"github.com/spf13/cobra"
)

const listHeight = 14

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("12"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)

type item string

func (i item) FilterValue() string { return "" }

type itemDelegate struct{}

func (d itemDelegate) Height() int                               { return 1 }
func (d itemDelegate) Spacing() int                              { return 0 }
func (d itemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
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

type model struct {
	list     list.Model
	items    []item
	choice   string
	quitting bool
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c":
			m.quitting = true
			return m, tea.Quit

		case "enter":
			i, ok := m.list.SelectedItem().(item)
			if ok {
				m.choice = string(i)
			}
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	if m.choice != "" {
		return quitTextStyle.Render(fmt.Sprintf("Selecting %s", m.choice))
	}
	if m.quitting {
		return quitTextStyle.Render("Exiting...")
	}
	return "\n" + m.list.View()
}

var useCMD = &cobra.Command{
	Use:   "use",
	Short: "Select kubeconfig to use",
	Run:   useFunc,
}

// init adds this command and his flags
func init() {
	rootCMD.AddCommand(useCMD)
	useCMD.Flags().StringP("name", "n", "", "name of the kubeconfig")
}

func useFunc(cmd *cobra.Command, args []string) {
	var name string

	retcode := 0
	defer func() { os.Exit(retcode) }()

	db, err := kv.Open(dbPath)
	if err != nil {
		log.Error().Str("error", err.Error()).Str("db", dbPath).Msg("error opening kv database")
		retcode = 1
		return
	}
	defer db.CloseDB()

	userHome, err := os.UserHomeDir()
	if err != nil {
		log.Fatal().Str("error", err.Error()).Msg("error getting the homeDir of the user")
	}
	kcPath := userHome + "/.kube/config"

	if cmd.Flags().Changed("name") {
		name, _ = cmd.Flags().GetString("name")
		err = useWithName(db, name, kcPath)
	} else {
		err = useWithoutName(db, kcPath)
	}
	if err != nil {
		log.Error().Str("error", err.Error()).Str("db", dbPath).Msg("error selecting kubeconfig to use")
		retcode = 1
		return
	}
	log.Debug().Str("command", "use").Str("key", name).Str("result", "successful").Send()
}

func useWithName(db *kv.DB, name, kubeconfigPath string) error {
	log.Debug().Str("command", "use").Str("with name", "true").Str("name", name).Str("database", dbPath).Send()

	kc, err := db.GetKubeconfig(name)
	if err != nil {
		return err
	}
	err = os.WriteFile(kubeconfigPath, []byte(kc.Content), 0o644)
	if err != nil {
		return err
	}
	return nil
}

func useWithoutName(db *kv.DB, kubeconfigPath string) error {
	log.Debug().Str("command", "use").Str("with name", "false").Str("database", dbPath).Send()

	var listItems []list.Item

	items, err := db.GetAll()
	if err != nil {
		return err
	}
	for _, kc := range items {
		listItems = append(listItems, item(kc.Name))
	}

	const defaultWidth = 20

	l := list.NewModel(listItems, itemDelegate{}, defaultWidth, listHeight)
	l.Title = "Select kubeconfig:"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	m := model{list: l}
	if err := tea.NewProgram(m).Start(); err != nil {
		return fmt.Errorf("error creating TUI: %s", err)
	}

	fmt.Printf("%s\n", m.choice)

	return nil
}
