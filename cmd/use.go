package cmd

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ebarped/kubeswap/internal/tui"
	"github.com/ebarped/kubeswap/pkg/kv"
	"github.com/spf13/cobra"
)

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
	defer db.Close()

	if db.IsEmpty() {
		log.Info().Msg("there are no kubeconfigs in the database. Exiting...")
		retcode = 0
		return
	}

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
	log.Debug().Str("command", "use").Str("result", "successful").Send()
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
		listItems = append(listItems, tui.Item(kc.Name))
	}

	const defaultWidth = 30

	l := list.NewModel(listItems, tui.ItemDelegate{}, defaultWidth, tui.ListHeight)
	l.Title = "Select kubeconfig:"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = tui.TitleStyle
	l.Styles.PaginationStyle = tui.PaginationStyle
	l.Styles.HelpStyle = tui.HelpStyle

	// we create a new model
	// it has a list of items and a channel,
	// so bubbletea can send the selected item outside its runtime
	m := tui.NewModel(l)

	// new program will take the model, and call Init,
	// then Update and then View, and alternate between
	// these 2 when an event (tea.Msg) is triggered (when something happens)
	err = tea.NewProgram(m).Start()
	if err != nil {
		return fmt.Errorf("error creating TUI: %s", err)
	}

	// Print out the final choice.
	if m.Choice != "" {
		log.Debug().Str("command", "use").Str("with name", "false").Str("database", dbPath).Str("TUI - item selected", m.Choice).Send()
		err = useWithName(db, m.Choice, kubeconfigPath)
		if err != nil {
			return err
		}
	}

	return nil
}
