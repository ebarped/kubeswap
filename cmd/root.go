package cmd

import (
	"fmt"
	"io"
	"io/ioutil"
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
	kcRootDir string
	dbPath    string // path to the db file
)

// variable shared at package level (used by subcommands)
var log *zerolog.Logger

// adds all the flags of the root command
// persistentFlags are the ones that are common to all subcommands
func init() {
	userHome, err := os.UserHomeDir()
	if err != nil {
		log.Fatal().Str("error", err.Error()).Msg("error getting the homeDir of the user")
	}
	dbPath = userHome + "/.kube/kubeswap.db"

	rootCMD.PersistentFlags().StringVar(&logLevel, "loglevel", "info", "loglevel (info/debug)")
	rootCMD.PersistentFlags().StringVar(&dbPath, "db", dbPath, "db file path")

	rootCMD.AddCommand(completionCmd)
}

// execute common initial steps
func initConfig(cmd *cobra.Command, args []string) {
	log = logger.New(logLevel)
	log.Debug().Msgf("loglevel set to %s", log.GetLevel().String())
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

	// we have a name arg, use that as kubeconfig
	if len(args) == 1 {
		name := args[0]
		path := kcRootDir + name
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

	files, err := ioutil.ReadDir(kcRootDir)
	if err != nil {
		log.Error().Msg(err.Error())
	}

	var listItems []list.Item

	for _, f := range files {
		log.Debug().Msgf("loading kubeconfig %s from %s", f.Name(), kcRootDir+f.Name())
		kc, err := kubeconfig.New(f.Name(), kcRootDir+f.Name())
		if err != nil {
			log.Error().Str("file", f.Name()).Msg(err.Error())
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

	// This is where we'll listen for the choice the user makes in the Bubble
	// Tea program.
	result := make(chan string, 1)

	// we create a new model
	// it has a list of items and a channel,
	// so bubbletea can send the selected item outside its runtime
	m := tui.NewModel(l, result)

	// new program will take the model, and call Init,
	// then Update and then View, and alternate between
	// these 2 when an event (tea.Msg) is triggered (when something happens)
	err = tea.NewProgram(m).Start()
	if err != nil {
		log.Error().Str("error", err.Error()).Msg("error creating bubbleteam program")
		retcode = 1
		return
	}

	// Print out the final choice.
	choice := <-result
	if choice != "" {
		kc, err := kubeconfig.New(choice, kcRootDir+choice)
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
