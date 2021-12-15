package cmd

import (
	"fmt"
	"os"

	"github.com/ebarped/kubeswap/pkg/logger"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

// variables to store the rootCMD flags
var logLevel string

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
}

// execute common initial steps
func initConfig(cmd *cobra.Command, args []string) {
	log = logger.New(logLevel)
	log.Debug().Msgf("loglevel set to %s", log.GetLevel().String())
}

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
