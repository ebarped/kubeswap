package cmd

import (
	"fmt"
	"os"

	"github.com/ebarped/kubeswap/pkg/kv"
	"github.com/spf13/cobra"
)

var printCMD = &cobra.Command{
	Use:   "print -n <name>",
	Short: "Prints the content of the kubeconfig referenced by <name>",
	Run:   printFunc,
}

// init adds this command and his flags
func init() {
	rootCMD.AddCommand(printCMD)
	printCMD.Flags().StringP("name", "n", "", "name of the kubeconfig")
	printCMD.MarkFlagRequired("name")
}

func printFunc(cmd *cobra.Command, args []string) {
	name, _ := cmd.Flags().GetString("name")

	retcode := 0
	defer func() { os.Exit(retcode) }()

	log.Debug().Str("command", "print").Str("database", dbPath).Send()

	db, err := kv.Open(dbPath)
	if err != nil {
		log.Error().Str("error", err.Error()).Msg("error opening kv database")
	}
	defer db.CloseDB()

	kc, err := db.GetKubeconfig(name)
	if err != nil {
		log.Error().Str("error", err.Error()).Msgf("error getting kubeconfig from key")
		retcode = 1
		return
	}
	fmt.Println(kc.Content)
	log.Debug().Str("command", "print").Str("key", kc.Name).Str("result", "successful").Send()
}
