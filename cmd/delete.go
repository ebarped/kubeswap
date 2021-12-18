package cmd

import (
	"os"

	"github.com/ebarped/kubeswap/pkg/kv"
	"github.com/spf13/cobra"
)

var deleteCMD = &cobra.Command{
	Use:   "delete -n <name>",
	Short: "Deletes a kubeconfig from the database",
	Run:   deleteFunc,
}

// init adds this command and his flags
func init() {
	rootCMD.AddCommand(deleteCMD)
	deleteCMD.Flags().StringP("name", "n", "", "name of the kubeconfig")
	deleteCMD.MarkFlagRequired("name")
}

func deleteFunc(cmd *cobra.Command, args []string) {
	name, _ := cmd.Flags().GetString("name")

	retcode := 0
	defer func() { os.Exit(retcode) }()

	log.Debug().Str("command", "delete").Str("name", name).Str("database", dbPath).Send()

	db, err := kv.Open(dbPath)
	if err != nil {
		log.Error().Str("error", err.Error()).Str("db", dbPath).Msg("error opening kv database")
		retcode = 1
		return
	}
	defer db.CloseDB()

	log.Debug().Str("action", "delete kubeconfig from db").Str("key", name).Send()

	err = db.DeleteKubeconfig(name)
	if err != nil {
		log.Error().Str("error", err.Error()).Msg("error putting new key-value to the database")
		retcode = 1
		return
	}
	log.Debug().Str("command", "delete").Str("key", name).Str("result", "successful").Send()
}
