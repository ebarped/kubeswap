package cmd

import (
	"fmt"
	"os"

	"github.com/akrylysov/pogreb"
	"github.com/ebarped/kubeswap/pkg/logger"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

// variables to store the rootCMD flags
var (
	logLevel   string
	kubeconfig string
	dbPath     string
)

// this variables are shared on package level
var (
	log *zerolog.Logger
	DB  *pogreb.DB
)

var rootCMD = &cobra.Command{
	Use:              "kubeswap",
	Short:            "tool to manage multiple kubeconfig files and change between clusters easily",
	PersistentPreRun: initConfig,
}

// Execute adds all child commands to the root command, and sets flags
func Execute() {
	if err := rootCMD.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// adds all the flags of the root command
// persistentFlags are the ones that are common to all subcommands
func init() {
	userHome, err := os.UserHomeDir()
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
	dbPath = userHome + "/.kube/kubeswap.db"

	rootCMD.PersistentFlags().StringVar(&logLevel, "log-level", "info", "loglevel (info/debug)")
	rootCMD.PersistentFlags().StringVar(&kubeconfig, "kubeconfig", "", "kubeconfig path")
	rootCMD.PersistentFlags().StringVar(&dbPath, "db", dbPath, "db file path")
}

// execute common initial steps:
// - create package level logger and set the loglevel
// - set the package level var userHome (to use in the init() func)
// - open the db file
func initConfig(cmd *cobra.Command, args []string) {
	log = logger.New(logLevel)
	log.Info().Msgf("loglevel set to %s", log.GetLevel().String())
	openDB()
}

// openDB opens the db file
func openDB() {
	userHome, err := os.UserHomeDir()
	if err != nil {
		log.Fatal().Msg(err.Error())
	}

	dbPath := userHome + "/.kube/kubeswap.db"

	DB, err = pogreb.Open(dbPath, nil)
	if err != nil {
		log.Fatal().Msg(err.Error())
		return
	}
}

func CloseDB() {
	err := DB.Close()
	if err != nil {
		log.Fatal().Msg(err.Error())
		return
	}
}
