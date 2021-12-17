package cmd

import (
	"fmt"
	"os"

	"github.com/ebarped/kubeswap/pkg/kv"
	"github.com/spf13/cobra"
)

var printallCMD = &cobra.Command{
	Use:   "printall",
	Short: "Prints the content of all the kubeconfigs from the db",
	Run:   printallFunc,
}

// init adds this command and his flags
func init() {
	rootCMD.AddCommand(printallCMD)
}

func printallFunc(cmd *cobra.Command, args []string) {
	retcode := 0
	defer func() { os.Exit(retcode) }()

	log.Debug().Str("command", "printall").Str("database", dbPath).Send()

	db, err := kv.Open(dbPath)
	if err != nil {
		log.Error().Str("error", err.Error()).Msg("error opening kv database")
	}
	defer db.CloseDB()

	items, err := db.GetAll()
	if err != nil {
		log.Error().Str("error", err.Error()).Msg("error listing items from database")
		retcode = 1
		return
	}
	for _, kc := range items {
		fmt.Println(kc.Content)
	}

	log.Debug().Str("command", "printall").Str("result", "successful").Send()
}
