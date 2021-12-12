package cmd

import (
	"fmt"
	"os"

	"github.com/ebarped/kubeswap/pkg/logger"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

var logLevel string

var log *zerolog.Logger

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
	rootCMD.PersistentFlags().StringVar(&logLevel, "log-level", "info", "loglevel (info/debug")
}

// execute common initial steps:
// - create global logger and set the loglevel
func initConfig(cmd *cobra.Command, args []string) {
	log = logger.New(logLevel)
	log.Info().Msgf("loglevel set to %s", log.GetLevel().String())
}

// This function is executed on startup to check
// if the database file accesible and readable/writable
// if any of the conditions are failed, print message to inform the user to execute setup command of change --db flag
func checkState() {
	fmt.Println("TODO")
}
