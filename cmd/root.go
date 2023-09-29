package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ebarped/kubeswap/internal/tui"
	"github.com/ebarped/kubeswap/pkg/kubeconfig"
	"github.com/ebarped/kubeswap/pkg/logger"
	"github.com/pterm/pterm"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

// variables to store the rootCMD flags
var (
	logLevel  string
	kcRootDir string // rootDir of the default kubeconfig location ($HOME/.kube/)
	dbPath    string // path to the db file
)

// variable shared at package level (used by subcommands)
var log *zerolog.Logger

// adds all the flags of the root command
// persistentFlags are the ones that are common to all subcommands
func init() {
	rootCMD.PersistentFlags().StringVar(&logLevel, "loglevel", "info", "loglevel (info/debug)")
	rootCMD.PersistentFlags().StringVar(&dbPath, "db", "$HOME/.kube/kubeswap.db", "db path")
}

// execute common initial steps
func initConfig(cmd *cobra.Command, args []string) {
	// create a comon logger
	log = logger.New(logLevel)
	log.Debug().Msgf("loglevel set to %s", log.GetLevel().String())

	// expand the dbPath var
	dbPath = os.ExpandEnv(dbPath)
}

var rootCMD = &cobra.Command{
	Use:              "kubeswap",
	Short:            "tool to manage multiple kubeconfig files and change between clusters easily",
	Long:             printLogo(),
	PersistentPreRun: initConfig,
	Run:              rootFunc,
	Args:             validateArgs,
}

func validateArgs(cmd *cobra.Command, args []string) error {
	// probably i can do better :rofl:
	if len(args) != 0 && len(args) != 1 {
		return fmt.Errorf("invalid number of args")
	}
	return nil
}

func rootFunc(cmd *cobra.Command, args []string) {
	retcode := 0
	defer func() { os.Exit(retcode) }()

	userHome, err := os.UserHomeDir()
	if err != nil {
		log.Fatal().Str("error", err.Error()).Msg("error getting the homeDir of the user")
	}
	kcRootDir = userHome + "/.kube/"

	// we have a name arg, use that as filename to select the kubeconfig
	if len(args) == 1 {
		name := args[0]
		path := kcRootDir + name

		// "deselect" kubeconfig
		if name == "none" {
			log.Debug().Str("name", name).Str("path", path).Msgf("removing default kubeconfig")
			err := deleteDefaultKubeconfig()
			if err != nil {
				log.Error().Str("name", "none").Str("path", path).Str("error", err.Error()).Msg("error setting kubeconfig to none...")
				retcode = 1
			}
			retcode = 0
			return
		}

		log.Debug().Str("name", name).Str("path", path).Msgf("loading kubeconfig...")
		kc, err := kubeconfig.New(name, path)
		if err != nil {
			log.Error().Str("name", name).Str("path", path).Str("error", err.Error()).Msg("error loading kubeconfig")
			retcode = 1
			return
		}
		err = useKubeconfig(kc)
		if err != nil {
			log.Error().Str("name", name).Str("path", path).Str("error", err.Error()).Msg("error selecting kubeconfig")
			retcode = 1
			return
		}
		log.Debug().Str("name", name).Str("path", path).Msg("kubeconfig successfully loaded")
		return
	}
	// we dont have the filename arg, so we scan the $HOME/.kube/ directory
	files, err := os.ReadDir(kcRootDir)
	if err != nil {
		log.Fatal().Msg(err.Error())
	}

	// check if we have any kubeconfig to list
	if len(files) == 0 {
		log.Info().Str("path", kcRootDir).Msg("Seems that you dont have any kubeconfig files ...")
		os.Exit(0)
	}

	var listItems []list.Item

	for _, f := range files {
		log.Debug().Str("file", f.Name()).Str("path", kcRootDir+f.Name()).Msg("loading kubeconfig...")
		// skip the default kubeconfig
		if f.Name() == "config" {
			log.Debug().Str("file", f.Name()).Msg("skipping default kubeconfig")
			continue
		}
		kc, err := kubeconfig.New(f.Name(), kcRootDir+f.Name())
		if err != nil {
			log.Debug().Str("file", f.Name()).Msg("not a valid kubeconfig")
			continue
		}
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

	// newProgram will take the model, and call Init,
	// then Update and then View, and alternate between
	// these 2 when an event (tea.Msg) is triggered (when something happens)
	err = tea.NewProgram(m).Start()
	if err != nil {
		log.Error().Str("error", err.Error()).Msg("error creating bubbleteam program")
		retcode = 1
		return
	}

	// Once the BubbleTeam runtime is done, we receive here the choice
	if m.Choice == "" {
		log.Debug().Msg("exit without selecting any item")
		return
	}

	kc, err := kubeconfig.New(m.Choice, kcRootDir+m.Choice)
	if err != nil {
		log.Error().Str("file", kc.Name).Msg(err.Error())
		retcode = 1
		return
	}

	err = useKubeconfig(kc)
	if err != nil {
		log.Error().Str("file", kc.Name).Msg(err.Error())
		retcode = 1
		return
	}
}

func deleteDefaultKubeconfig() error {
	path := kcRootDir + "config"

	_, err := os.Stat(path)
	if err != nil {
		// .kube/config does not exist, just exit
		return nil
	}
	err = os.Remove(path)
	if err != nil {
		return err
	}
	return nil
}

func useKubeconfig(kc *kubeconfig.Kubeconfig) error {
	_, err := copy(kc.Path, kcRootDir+"config")
	if err != nil {
		return err
	}
	return nil
}

// Execute adds all child commands to the root command, and sets flags
func Execute() {
	if err := rootCMD.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func printLogo() string {
	logo := pterm.FgBlue.Sprint(`
  ██   ██ ██    ██ ██████  ███████ ███████ ██     ██  █████  ██████  
  ██  ██  ██    ██ ██   ██ ██      ██      ██     ██ ██   ██ ██   ██ 
  █████   ██    ██ ██████  █████   ███████ ██  █  ██ ███████ ██████  
  ██  ██  ██    ██ ██   ██ ██           ██ ██ ███ ██ ██   ██ ██      
  ██   ██  ██████  ██████  ███████ ███████  ███ ███  ██   ██ ██
`)

	subtext := pterm.FgLightCyan.Sprintf("Manage your kubeconfig files easily")
	return fmt.Sprintf(`
%s
%s
`, logo, subtext)
}

func copy(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}
